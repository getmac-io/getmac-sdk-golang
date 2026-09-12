package getmac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VirtualMachine struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Region       string    `json:"region"`
	Type         string    `json:"type"`
	StatusReason string    `json:"status_reason"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateVirtualMachineRequest describes a virtual machine to create.
//
// Image is either a GetMac runner label, the names GitHub Actions uses in
// runs-on (for example "getmac" or "getmac-tahoe"), or an image slug (for
// example "macos-tahoe"). A label also selects the machine type, so Type can be
// left empty to use it; an image slug needs Type.
type CreateVirtualMachineRequest struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Region string `json:"region"`
	Type   string `json:"type"`
}

type ListVirtualMachinesResponse struct {
	*ListResponse
	VirtualMachines []*VirtualMachine `json:"instances"`
}

type virtualMachinesService struct {
	client *Client
}

func (s *virtualMachinesService) Create(
	ctx context.Context, projectID string, req *CreateVirtualMachineRequest) (*http.Response, *VirtualMachine, error) {
	resp, err := s.client.doRequest(ctx, http.MethodPost, fmt.Sprintf("/instances?project_id=%s", projectID), req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, nil, newAPIError(resp)
	}

	var virtualMachine VirtualMachine
	if err := json.NewDecoder(resp.Body).Decode(&virtualMachine); err != nil {
		return resp, nil, err
	}

	return resp, &virtualMachine, nil
}

func (s *virtualMachinesService) Get(
	ctx context.Context, projectID string, id string) (*http.Response, *VirtualMachine, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/instances/%s?project_id=%s", id, projectID), nil)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, nil, newAPIError(resp)
	}

	var virtualMachine VirtualMachine
	if err := json.NewDecoder(resp.Body).Decode(&virtualMachine); err != nil {
		return resp, nil, err
	}

	return resp, &virtualMachine, nil
}

// GetByName returns the virtual machine with the given name. The error wraps
// ErrNotFound when no virtual machine has that name.
func (s *virtualMachinesService) GetByName(
	ctx context.Context, projectID string, name string) (*http.Response, *VirtualMachine, error) {
	resp, vms, err := s.List(ctx, projectID)
	if err != nil {
		return resp, nil, err
	}

	for _, vm := range vms {
		if vm.Name != name {
			continue
		}

		return resp, vm, nil
	}

	return resp, nil, fmt.Errorf("virtual machine with name %s %w", name, ErrNotFound)
}

func (s *virtualMachinesService) Delete(
	ctx context.Context, projectID string, id string) (*http.Response, error) {
	resp, err := s.client.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/instances/%s?project_id=%s", id, projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, newAPIError(resp)
	}

	return resp, nil
}

func (s *virtualMachinesService) List(
	ctx context.Context, projectID string) (*http.Response, []*VirtualMachine, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/instances?project_id=%s", projectID), nil)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, nil, newAPIError(resp)
	}

	var listResp ListVirtualMachinesResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, nil, err
	}

	return resp, listResp.VirtualMachines, nil
}

func (s *virtualMachinesService) Start(
	ctx context.Context, projectID string, id string) (*http.Response, error) {
	resp, err := s.client.doRequest(ctx, http.MethodPost, fmt.Sprintf("/instances/%s/start?project_id=%s", id, projectID), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, newAPIError(resp)
	}

	return resp, nil
}

func (s *virtualMachinesService) Stop(
	ctx context.Context, projectID string, id string) (*http.Response, error) {
	resp, err := s.client.doRequest(ctx, http.MethodPost, fmt.Sprintf("/instances/%s/stop?project_id=%s", id, projectID), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp, newAPIError(resp)
	}

	return resp, nil
}
