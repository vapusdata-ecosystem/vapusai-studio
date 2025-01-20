package utils

import (
	models "github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func AIMNM2Pb(obj *models.AIModelNode) *mpb.AIModelNode {
	return obj.ConvertToPb()
}

func AIMNPb2Obj(obj *mpb.AIModelNode) *models.AIModelNode {
	return (&models.AIModelNode{}).ConvertFromPb(obj)
}

func AIMNListM2Pb(obj []*models.AIModelNode) []*mpb.AIModelNode {
	var res []*mpb.AIModelNode
	for _, v := range obj {
		res = append(res, AIMNM2Pb(v))
	}
	return res
}

func AIMNListPb2Obj(obj []*mpb.AIModelNode) []*models.AIModelNode {
	var res []*models.AIModelNode
	for _, v := range obj {
		res = append(res, AIMNPb2Obj(v))
	}
	return res
}

func AIPRObj2Pb(obj []*models.AIModelPrompt) []*mpb.AIModelPrompt {
	var res []*mpb.AIModelPrompt
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

func AIPRPb2Obj(obj []*mpb.AIModelPrompt) []*models.AIModelPrompt {
	var res []*models.AIModelPrompt
	for _, v := range obj {
		res = append(res, (&models.AIModelPrompt{}).ConvertFromPb(v))
	}
	return res
}

func AIAGPb2Obj(obj []*models.VapusAIAgent) []*mpb.VapusAIAgent {
	var res []*mpb.VapusAIAgent
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

func AIGDPb2Obj(obj []*mpb.AIGuardrails) []*models.AIGuardrails {
	var res []*models.AIGuardrails
	for _, v := range obj {
		res = append(res, (&models.AIGuardrails{}).ConvertFromPb(v))
	}
	return res
}

func AIGDObjToPb(obj []*models.AIGuardrails) []*mpb.AIGuardrails {
	var res []*mpb.AIGuardrails
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

// Utils func to convert proto domain  creation request to local obj
func DmNodeToObj(request *pb.OrganizationManagerRequest) *models.Organization {
	return (&models.Organization{}).ConvertFromPb(request.GetSpec())
}

// Utils func to convert local domain  object to proto object
func DmNToPb(objList []*models.Organization) []*mpb.Organization {
	if objList == nil {
		return nil
	}
	var pbList []*mpb.Organization
	for _, obj := range objList {
		pbList = append(pbList, obj.ConvertToPb())
	}
	return pbList
}

func DmUToPb(users []*models.Users, domain string) []*mpb.User {
	if users == nil {
		return nil
	}
	var pbList []*mpb.User
	for _, obj := range users {
		pbList = append(pbList, obj.ConvertToPb(domain))
	}
	return pbList
}

func DPPLObj2Pb(dws []*models.Plugin) []*mpb.Plugin {
	if dws == nil {
		return nil
	}
	var dpList []*mpb.Plugin
	for _, dp := range dws {
		dpList = append(dpList, dp.ConvertToPb())
	}
	return dpList
}

func DmUPbToObj(users *mpb.User) *models.Users {
	if users == nil {
		return nil
	}
	return (&models.Users{}).ConvertFromPb(users)
}
