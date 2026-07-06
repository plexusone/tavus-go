package tavus

import (
	"context"
	"fmt"

	"github.com/plexusone/tavus-go/api"
)

// ConversationsService handles real-time conversation sessions.
type ConversationsService struct {
	client *api.Client
}

// Create creates a new conversation session.
func (s *ConversationsService) Create(ctx context.Context, req *api.CreateConversationRequest) (*api.CreateConversationResponse, error) {
	res, err := s.client.CreateConversation(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreateConversationResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all conversations.
func (s *ConversationsService) List(ctx context.Context, params api.ListConversationsParams) (*api.ListConversationsResponse, error) {
	res, err := s.client.ListConversations(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListConversationsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific conversation.
func (s *ConversationsService) Get(ctx context.Context, conversationID string) (*api.Conversation, error) {
	res, err := s.client.GetConversation(ctx, api.GetConversationParams{ConversationID: conversationID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Conversation); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a conversation.
func (s *ConversationsService) Delete(ctx context.Context, conversationID string) error {
	_, err := s.client.DeleteConversation(ctx, api.DeleteConversationParams{ConversationID: conversationID})
	return err
}

// End gracefully ends an active conversation.
func (s *ConversationsService) End(ctx context.Context, conversationID string) error {
	_, err := s.client.EndConversation(ctx, api.EndConversationParams{ConversationID: conversationID})
	return err
}

// PalsService handles PAL (Personality AI Layer) management.
type PalsService struct {
	client *api.Client
}

// Create creates a new PAL.
func (s *PalsService) Create(ctx context.Context, req *api.CreatePalRequest) (*api.CreatePalResponse, error) {
	res, err := s.client.CreatePal(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreatePalResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all PALs.
func (s *PalsService) List(ctx context.Context, params api.ListPalsParams) (*api.ListPalsResponse, error) {
	res, err := s.client.ListPals(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListPalsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific PAL.
func (s *PalsService) Get(ctx context.Context, palID string) (*api.Pal, error) {
	res, err := s.client.GetPal(ctx, api.GetPalParams{PalID: palID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Pal); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates a PAL.
func (s *PalsService) Update(ctx context.Context, palID string, req *api.UpdatePalRequest) (*api.Pal, error) {
	res, err := s.client.UpdatePal(ctx, req, api.UpdatePalParams{PalID: palID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Pal); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a PAL.
func (s *PalsService) Delete(ctx context.Context, palID string) error {
	_, err := s.client.DeletePal(ctx, api.DeletePalParams{PalID: palID})
	return err
}

// AttachTool attaches a tool to a PAL.
func (s *PalsService) AttachTool(ctx context.Context, palID, toolID string) error {
	_, err := s.client.AttachToolToPal(ctx, &api.AttachToolToPalReq{ToolID: toolID}, api.AttachToolToPalParams{PalID: palID})
	return err
}

// FacesService handles face (visual identity) management.
type FacesService struct {
	client *api.Client
}

// Create starts training a new face.
func (s *FacesService) Create(ctx context.Context, req *api.CreateFaceRequest) (*api.CreateFaceResponse, error) {
	res, err := s.client.CreateFace(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreateFaceResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all faces.
func (s *FacesService) List(ctx context.Context, params api.ListFacesParams) (*api.ListFacesResponse, error) {
	res, err := s.client.ListFaces(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListFacesResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific face.
func (s *FacesService) Get(ctx context.Context, faceID string) (*api.Face, error) {
	res, err := s.client.GetFace(ctx, api.GetFaceParams{FaceID: faceID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Face); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a face.
func (s *FacesService) Delete(ctx context.Context, faceID string) error {
	_, err := s.client.DeleteFace(ctx, api.DeleteFaceParams{FaceID: faceID})
	return err
}

// Rename renames a face.
func (s *FacesService) Rename(ctx context.Context, faceID, name string) error {
	_, err := s.client.RenameFace(ctx, &api.RenameFaceReq{FaceName: name}, api.RenameFaceParams{FaceID: faceID})
	return err
}

// VoicesService handles stock voice listing.
type VoicesService struct {
	client *api.Client
}

// List returns available stock voices.
func (s *VoicesService) List(ctx context.Context, params api.ListVoicesParams) (*api.ListVoicesResponse, error) {
	res, err := s.client.ListVoices(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListVoicesResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// VideosService handles async video generation.
type VideosService struct {
	client *api.Client
}

// Create starts generating a new video.
func (s *VideosService) Create(ctx context.Context, req *api.CreateVideoRequest) (*api.CreateVideoResponse, error) {
	res, err := s.client.CreateVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreateVideoResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all videos.
func (s *VideosService) List(ctx context.Context, params api.ListVideosParams) (*api.ListVideosResponse, error) {
	res, err := s.client.ListVideos(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListVideosResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific video.
func (s *VideosService) Get(ctx context.Context, videoID string) (*api.Video, error) {
	res, err := s.client.GetVideo(ctx, api.GetVideoParams{VideoID: videoID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Video); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, videoID string) error {
	_, err := s.client.DeleteVideo(ctx, api.DeleteVideoParams{VideoID: videoID})
	return err
}

// Rename renames a video.
func (s *VideosService) Rename(ctx context.Context, videoID, name string) error {
	_, err := s.client.RenameVideo(ctx, &api.RenameVideoReq{VideoName: name}, api.RenameVideoParams{VideoID: videoID})
	return err
}

// ToolsService handles function calling tool management.
type ToolsService struct {
	client *api.Client
}

// Create creates a new tool.
func (s *ToolsService) Create(ctx context.Context, req *api.CreateToolRequest) (*api.Tool, error) {
	res, err := s.client.CreateTool(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Tool); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all tools.
func (s *ToolsService) List(ctx context.Context, params api.ListToolsParams) (*api.ListToolsResponse, error) {
	res, err := s.client.ListTools(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListToolsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific tool.
func (s *ToolsService) Get(ctx context.Context, toolID string) (*api.Tool, error) {
	res, err := s.client.GetTool(ctx, api.GetToolParams{ToolID: toolID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Tool); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates a tool.
func (s *ToolsService) Update(ctx context.Context, toolID string, req *api.UpdateToolRequest) (*api.Tool, error) {
	res, err := s.client.UpdateTool(ctx, req, api.UpdateToolParams{ToolID: toolID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Tool); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a tool.
func (s *ToolsService) Delete(ctx context.Context, toolID string) error {
	_, err := s.client.DeleteTool(ctx, api.DeleteToolParams{ToolID: toolID})
	return err
}

// GuardrailsService handles behavioral boundary management.
type GuardrailsService struct {
	client *api.Client
}

// Create creates a new guardrail.
func (s *GuardrailsService) Create(ctx context.Context, req *api.CreateGuardrailRequest) (*api.Guardrail, error) {
	res, err := s.client.CreateGuardrail(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Guardrail); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all guardrails.
func (s *GuardrailsService) List(ctx context.Context, params api.ListGuardrailsParams) (*api.ListGuardrailsResponse, error) {
	res, err := s.client.ListGuardrails(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListGuardrailsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific guardrail.
func (s *GuardrailsService) Get(ctx context.Context, guardrailID string) (*api.Guardrail, error) {
	res, err := s.client.GetGuardrail(ctx, api.GetGuardrailParams{GuardrailID: guardrailID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Guardrail); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates a guardrail.
func (s *GuardrailsService) Update(ctx context.Context, guardrailID string, req *api.UpdateGuardrailRequest) (*api.Guardrail, error) {
	res, err := s.client.UpdateGuardrail(ctx, req, api.UpdateGuardrailParams{GuardrailID: guardrailID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Guardrail); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a guardrail.
func (s *GuardrailsService) Delete(ctx context.Context, guardrailID string) error {
	_, err := s.client.DeleteGuardrail(ctx, api.DeleteGuardrailParams{GuardrailID: guardrailID})
	return err
}

// ObjectivesService handles conversation goal management.
type ObjectivesService struct {
	client *api.Client
}

// Create creates new objectives.
func (s *ObjectivesService) Create(ctx context.Context, req *api.CreateObjectivesRequest) (*api.CreateObjectivesResponse, error) {
	res, err := s.client.CreateObjectives(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreateObjectivesResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all objectives.
func (s *ObjectivesService) List(ctx context.Context, params api.ListObjectivesParams) (*api.ListObjectivesResponse, error) {
	res, err := s.client.ListObjectives(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListObjectivesResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns specific objectives.
func (s *ObjectivesService) Get(ctx context.Context, objectivesID string) (*api.Objectives, error) {
	res, err := s.client.GetObjectives(ctx, api.GetObjectivesParams{ObjectivesID: objectivesID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Objectives); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates objectives.
func (s *ObjectivesService) Update(ctx context.Context, objectivesID string, req *api.UpdateObjectivesRequest) (*api.Objectives, error) {
	res, err := s.client.UpdateObjectives(ctx, req, api.UpdateObjectivesParams{ObjectivesID: objectivesID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Objectives); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes objectives.
func (s *ObjectivesService) Delete(ctx context.Context, objectivesID string) error {
	_, err := s.client.DeleteObjectives(ctx, api.DeleteObjectivesParams{ObjectivesID: objectivesID})
	return err
}

// DocumentsService handles knowledge base document management.
type DocumentsService struct {
	client *api.Client
}

// Create uploads a new document.
func (s *DocumentsService) Create(ctx context.Context, req *api.CreateDocumentRequest) (*api.CreateDocumentResponse, error) {
	res, err := s.client.CreateDocument(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.CreateDocumentResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all documents.
func (s *DocumentsService) List(ctx context.Context, params api.ListDocumentsParams) (*api.ListDocumentsResponse, error) {
	res, err := s.client.ListDocuments(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListDocumentsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific document.
func (s *DocumentsService) Get(ctx context.Context, documentID string) (*api.Document, error) {
	res, err := s.client.GetDocument(ctx, api.GetDocumentParams{DocumentID: documentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Document); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates document metadata.
func (s *DocumentsService) Update(ctx context.Context, documentID string, req *api.UpdateDocumentRequest) (*api.Document, error) {
	res, err := s.client.UpdateDocument(ctx, req, api.UpdateDocumentParams{DocumentID: documentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Document); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a document.
func (s *DocumentsService) Delete(ctx context.Context, documentID string) error {
	_, err := s.client.DeleteDocument(ctx, api.DeleteDocumentParams{DocumentID: documentID})
	return err
}

// DeploymentsService handles distribution channel management.
type DeploymentsService struct {
	client *api.Client
}

// Create creates a new deployment.
func (s *DeploymentsService) Create(ctx context.Context, req *api.CreateDeploymentRequest) (*api.Deployment, error) {
	res, err := s.client.CreateDeployment(ctx, req)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Deployment); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// List returns all deployments.
func (s *DeploymentsService) List(ctx context.Context, params api.ListDeploymentsParams) (*api.ListDeploymentsResponse, error) {
	res, err := s.client.ListDeployments(ctx, params)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ListDeploymentsResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Get returns a specific deployment.
func (s *DeploymentsService) Get(ctx context.Context, deploymentID string) (*api.Deployment, error) {
	res, err := s.client.GetDeployment(ctx, api.GetDeploymentParams{DeploymentID: deploymentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Deployment); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Update updates a deployment.
func (s *DeploymentsService) Update(ctx context.Context, deploymentID string, req *api.UpdateDeploymentRequest) (*api.Deployment, error) {
	res, err := s.client.UpdateDeployment(ctx, req, api.UpdateDeploymentParams{DeploymentID: deploymentID})
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.Deployment); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// Delete deletes a deployment.
func (s *DeploymentsService) Delete(ctx context.Context, deploymentID string) error {
	_, err := s.client.DeleteDeployment(ctx, api.DeleteDeploymentParams{DeploymentID: deploymentID})
	return err
}
