package service

import (
	"context"
	"fmt"
	"github.com/cossim/coss-server/interface/live/api/dto"
	"github.com/cossim/coss-server/interface/live/api/model"
	msgconfig "github.com/cossim/coss-server/interface/msg/config"
	"github.com/cossim/coss-server/pkg/cache"
	"github.com/cossim/coss-server/pkg/code"
	groupgrpcv1 "github.com/cossim/coss-server/service/group/api/v1"
	relationgrpcv1 "github.com/cossim/coss-server/service/relation/api/v1"
	usergrpcv1 "github.com/cossim/coss-server/service/user/api/v1"
	"github.com/google/uuid"
	"github.com/livekit/protocol/livekit"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (s *Service) CreateGroupCall(ctx context.Context, uid string, gid uint32, member []string) (interface{}, error) {
	is, err := s.isEitherGroupInCall(ctx, gid)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, code.LiveErrAlreadyInCall
	}

	group, err := s.groupClient.GetGroupInfoByGid(ctx, &groupgrpcv1.GetGroupInfoRequest{Gid: gid})
	if err != nil {
		s.logger.Error("create group call failed", zap.Error(err))
		return nil, err
	}

	if group.Status != groupgrpcv1.GroupStatus_GROUP_STATUS_NORMAL {
		return nil, code.GroupErrGroupStatusNotAvailable
	}

	//rel, err := s.relGroupClient.GetGroupRelation(ctx, &relationgrpcv1.GetGroupRelationRequest{GroupId: gid, UserId: uid})
	//if err != nil {
	//	return nil, err
	//}
	//
	//if rel.MuteEndTime > 0 {
	//	return nil, code.GroupErrUserIsMuted
	//}

	type MemberInfo struct {
		UserID   string
		UserName string
		Token    string
	}

	rels, err := s.relGroupClient.GetBatchGroupRelation(ctx, &relationgrpcv1.GetBatchGroupRelationRequest{GroupId: gid, UserIds: member})
	if err != nil {
		return nil, err
	}

	if len(member) != len(rels.GroupRelationResponses) {
		return nil, code.RelationGroupErrNotInGroup
	}

	roomName := uuid.New().String()
	room, err := s.roomService.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    uint32(s.liveTimeout.Seconds()), // 空闲时间
		MaxParticipants: 256,
	})
	if err != nil {
		s.logger.Error("创建通话失败", zap.Error(err))
		return nil, err
	}

	redisRoom := &model.GroupRoomInfo{
		Room:            room.Name,
		GroupID:         gid,
		SenderID:        uid,
		MaxParticipants: 256,
		Participants:    map[string]*model.ActiveParticipant{},
	}

	// 提取未被禁言的成员信息
	participants := make([]MemberInfo, 0)
	participants = append(participants, MemberInfo{UserID: uid})
	for _, v := range rels.GroupRelationResponses {
		//if v.MuteEndTime == 0 {
		participants = append(participants, MemberInfo{UserID: v.UserId})
		//}
	}

	var senderToken string

	// 获取成员的名称和其他信息
	for i := range participants {
		memberUser, err := s.userClient.UserInfo(ctx, &usergrpcv1.UserInfoRequest{UserId: participants[i].UserID})
		if err != nil {
			s.logger.Error("获取用户信息失败", zap.Error(err))
			continue
		}

		if memberUser.Status != usergrpcv1.UserStatus_USER_STATUS_NORMAL {
			s.logger.Error("用户状态异常", zap.String(memberUser.UserId, memberUser.Status.String()))
			continue
		}

		//redisRoom.Participants[participants[i].UserID] = &model.ActiveParticipant{}
		participants[i].UserName = memberUser.NickName

		if participants[i].UserID == uid {
			senderToken, err = s.GetJoinToken(ctx, room.Name, "admin", memberUser.NickName)
			if err != nil {
				s.logger.Error("获取token失败", zap.Error(err))
				continue
			}
			participants[i].Token = senderToken
			continue
		}

		token, err := s.GetJoinToken(ctx, room.Name, "user", memberUser.NickName)
		if err != nil {
			s.logger.Error("获取token失败", zap.Error(err))
			continue
		}

		msg := msgconfig.WsMsg{Uid: participants[i].UserID, Event: msgconfig.GroupCallReqEvent, Data: map[string]interface{}{
			"url":          s.livekitServer,
			"token":        token,
			"room":         room.Name,
			"sender_id":    uid,
			"recipient_id": participants[i].UserID,
		}}
		if err = s.publishServiceMessage(ctx, msg); err != nil {
			s.logger.Error("发送消息失败", zap.Error(err))
		}
	}

	ToJSONString, err := redisRoom.ToJSONString()
	if err != nil {
		return nil, err
	}
	fmt.Println("ToJSONString => ", ToJSONString)
	if err = cache.SetKey(s.redisClient, liveGroupPrefix+room.Name+":"+strconv.FormatUint(uint64(gid), 10), ToJSONString, s.liveTimeout); err != nil {
		s.logger.Error("保存房间信息失败", zap.Error(err))
		return nil, err
	}

	return map[string]string{
		"url":       s.livekitServer,
		"token":     senderToken,
		"room_name": roomName,
		"room_id":   room.Sid,
	}, nil
}

func (s *Service) GroupJoinRoom(ctx context.Context, gid uint32, uid string) (interface{}, error) {
	_, err := s.relGroupClient.GetGroupRelation(ctx, &relationgrpcv1.GetGroupRelationRequest{GroupId: gid, UserId: uid})
	if err != nil {
		s.logger.Error("获取用户群组关系失败", zap.Error(err))
		return nil, err
	}

	roomInfo, key, err := s.getGroupRedisRoom(ctx, gid)
	if err != nil {
		s.logger.Error("获取群聊房间信息失败", zap.Error(err))
		return nil, code.LiveErrGetCallInfoFailed
	}

	if err = s.checkGroupRoom(ctx, roomInfo, uid, roomInfo.Room); err != nil {
		return nil, err
	}

	roomInfo.NumParticipants++
	roomInfo.Participants[uid] = &model.ActiveParticipant{
		Connected: true,
	}
	ToJSONString, err := roomInfo.ToJSONString()
	if err != nil {
		return nil, err
	}
	if err = cache.SetKey(s.redisClient, key, ToJSONString, 0); err != nil {
		s.logger.Error("更新房间信息失败", zap.Error(err))
		return nil, err
	}

	return nil, nil
}

func (s *Service) GroupShowRoom(ctx context.Context, gid uint32, uid string) (*dto.GroupShowResponse, error) {
	_, err := s.relGroupClient.GetGroupRelation(ctx, &relationgrpcv1.GetGroupRelationRequest{GroupId: gid, UserId: uid})
	if err != nil {
		s.logger.Error("获取用户群组关系失败", zap.Error(err))
		return nil, err
	}

	room, key, err := s.getGroupRedisRoom(ctx, gid)
	if err != nil {
		s.logger.Error("获取群聊房间信息失败", zap.Error(err))
		return nil, code.LiveErrGetCallInfoFailed
	}

	livekitRoom, err := s.getLivekitRoom(ctx, room.Room, func(ctx context.Context, room string) {
		s.deleteGroupRoom(ctx, room, key)
	})
	if err != nil {
		return nil, err
	}

	resp := &dto.GroupShowResponse{
		GroupID:         room.GroupID,
		NumParticipants: livekitRoom.NumParticipants,
		MaxParticipants: livekitRoom.MaxParticipants,
		StartAt:         livekitRoom.CreationTime,
		Duration:        0,
		Room:            livekitRoom.Name,
		Participant:     make([]*dto.ParticipantInfo, 0),
	}
	res, err := s.roomService.ListParticipants(ctx, &livekit.ListParticipantsRequest{
		Room: room.Room, // 房间名称
	})
	if err != nil {
		s.logger.Error("获取通话信息失败", zap.Error(err))
		return nil, code.LiveErrGetCallInfoFailed
	}

	for _, p := range res.Participants {
		//if p.Identity == userID {
		//	continue
		//}
		resp.Participant = append(resp.Participant, &dto.ParticipantInfo{
			Sid:         p.Sid,
			Identity:    p.Identity,
			State:       dto.ParticipantState(p.State),
			JoinedAt:    p.JoinedAt,
			Name:        p.Name,
			IsPublisher: p.IsPublisher,
		})
	}

	return resp, nil
}

func (s *Service) GroupRejectRoom(ctx context.Context, gid uint32, uid string) (interface{}, error) {
	_, err := s.relGroupClient.GetGroupRelation(ctx, &relationgrpcv1.GetGroupRelationRequest{GroupId: gid, UserId: uid})
	if err != nil {
		s.logger.Error("获取用户群组关系失败", zap.Error(err))
		return nil, err
	}

	roomInfo := &model.UserRoomInfo{}
	_, _, err = s.getGroupRedisRoom(ctx, gid)
	if err != nil {
		return nil, err
	}

	//_, err = s.deleteUserRoom(ctx, room)
	//if err != nil {
	//	return nil, err
	//}

	msg := msgconfig.WsMsg{Uid: roomInfo.SenderID, Event: msgconfig.GroupCallRejectEvent, Data: map[string]interface{}{
		"sender_id":    roomInfo.SenderID,
		"recipient_id": uid,
	}}
	if err = s.publishServiceMessage(ctx, msg); err != nil {
		s.logger.Error("发送消息失败", zap.Error(err))
		return nil, err
	}

	return nil, nil
}

func (s *Service) GroupLeaveRoom(ctx context.Context, gid uint32, uid string, force bool) (interface{}, error) {
	_, err := s.relGroupClient.GetGroupRelation(ctx, &relationgrpcv1.GetGroupRelationRequest{GroupId: gid, UserId: uid})
	if err != nil {
		s.logger.Error("获取用户群组关系失败", zap.Error(err))
		return nil, err
	}

	roomInfo, key, err := s.getGroupRedisRoom(ctx, gid)
	if err != nil {
		return nil, err
	}

	if roomInfo.SenderID != uid && force {
		return nil, code.Unauthorized
	}

	if force || roomInfo.NumParticipants-1 == 0 || roomInfo.NumParticipants == 0 {
		_, err = s.deleteGroupRoom(ctx, roomInfo.Room, key)
		if err != nil {
			return nil, err
		}
	} else {
		delete(roomInfo.Participants, uid)
		if roomInfo.NumParticipants > 0 {
			roomInfo.NumParticipants--
		}
		ToJSONString, err := roomInfo.ToJSONString()
		if err != nil {
			return nil, err
		}
		fmt.Println("111111111")
		if err = cache.SetKey(s.redisClient, key, ToJSONString, 0); err != nil {
			s.logger.Error("更新房间信息失败", zap.Error(err))
			return nil, err
		}
	}

	for k := range roomInfo.Participants {
		if k == uid {
			continue
		}
		msg := msgconfig.WsMsg{Uid: k, Event: msgconfig.GroupCallEndEvent, Data: map[string]interface{}{
			"sender_id":    roomInfo.SenderID,
			"recipient_id": uid,
		}}
		if err = s.publishServiceMessage(ctx, msg); err != nil {
			s.logger.Error("发送消息失败", zap.Error(err))
			return nil, err
		}
	}

	return nil, nil
}

func (s *Service) getGroupRedisRoom(ctx context.Context, gid uint32) (*model.GroupRoomInfo, string, error) {
	k1 := "*:" + strconv.FormatUint(uint64(gid), 10)
	s.logger.Info("getGroupRedisRoom", zap.String("key", k1))
	resp := &model.GroupRoomInfo{}
	k2, err := s.getRedisRoomWithPrefix(ctx, liveGroupPrefix, k1, resp)
	if err != nil {
		return nil, "", err
	}
	return resp, k2, nil
}

func (s *Service) isEitherGroupInCall(ctx context.Context, gid uint32) (bool, error) {
	pattern := liveGroupPrefix + "*:" + strconv.FormatUint(uint64(gid), 10)
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Service) deleteGroupRoom(ctx context.Context, room string, key string) (interface{}, error) {
	_, err := s.roomService.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
		Room: room, // 房间名称
	})
	if err != nil && !strings.Contains(err.Error(), "room not found") {
		s.logger.Error("error deleting room", zap.Error(err))
		return nil, code.LiveErrLeaveCallFailed
	}

	_, err = s.redisClient.Del(ctx, key).Result()
	if err != nil {
		s.logger.Error("删除房间信息失败", zap.Error(err))
	}

	return nil, nil
}

func (s *Service) checkGroupRoom(ctx context.Context, roomInfo *model.GroupRoomInfo, uid, roomName string) error {
	//if uid != roomInfo.RecipientID && uid != roomInfo.SenderID {
	//	return code.Unauthorized
	//}

	if roomInfo.NumParticipants+1 > roomInfo.MaxParticipants {
		return code.LiveErrMaxParticipantsExceeded
	}

	if _, ok := roomInfo.Participants[uid]; ok {
		return code.LiveErrAlreadyInCall
	}

	room, err := s.getLivekitRoom(ctx, roomName, nil)
	if err != nil {
		return err
	}

	if room.NumParticipants+1 > room.MaxParticipants {
		return code.LiveErrMaxParticipantsExceeded
	}

	return nil
}
