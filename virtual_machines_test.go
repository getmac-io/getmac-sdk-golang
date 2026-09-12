package getmac

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	return NewClient(WithBaseURL(srv.URL), WithToken("token"))
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func TestCreate_ReturnsAPIErrorWithMessage(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Image or label getmac-nope not found"})
	})

	_, vm, err := client.VirtualMachines().Create(context.Background(), "project", &CreateVirtualMachineRequest{
		Name: "vm", Image: "getmac-nope", Region: "eu-central-ltu-1"})
	if vm != nil {
		t.Fatalf("expected no virtual machine, got %+v", vm)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected an APIError, got %v", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest || apiErr.Message != "Image or label getmac-nope not found" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}

	if want := "unexpected status code: 400: Image or label getmac-nope not found"; err.Error() != want {
		t.Fatalf("expected error %q, got %q", want, err.Error())
	}

	if errors.Is(err, ErrNotFound) {
		t.Fatal("a 400 response must not match ErrNotFound")
	}
}

func TestAPIError_KeepsPlainTextBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("  bad gateway\n"))
	})

	_, err := client.VirtualMachines().Delete(context.Background(), "project", "vm-1")
	if want := "unexpected status code: 502: bad gateway"; err == nil || err.Error() != want {
		t.Fatalf("expected error %q, got %v", want, err)
	}
}

func TestAPIError_WithoutBody(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.VirtualMachines().Start(context.Background(), "project", "vm-1")
	if want := "unexpected status code: 500"; err == nil || err.Error() != want {
		t.Fatalf("expected error %q, got %v", want, err)
	}
}

func TestGet_NotFoundMatchesErrNotFound(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]string{"message": "Instance not found"})
	})

	resp, _, err := client.VirtualMachines().Get(context.Background(), "project", "vm-1")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected the 404 response to be returned, got %v", resp)
	}
}

func TestGetByName(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]any{
			"total":     1,
			"instances": []map[string]string{{"id": "vm-1", "name": "gitlab-job-1"}},
		})
	})

	t.Run("match", func(t *testing.T) {
		_, vm, err := client.VirtualMachines().GetByName(context.Background(), "project", "gitlab-job-1")
		if err != nil || vm == nil || vm.ID != "vm-1" {
			t.Fatalf("expected vm-1, got %+v, %v", vm, err)
		}
	})

	t.Run("no match", func(t *testing.T) {
		_, vm, err := client.VirtualMachines().GetByName(context.Background(), "project", "gitlab-job-2")
		if vm != nil {
			t.Fatalf("expected no virtual machine, got %+v", vm)
		}

		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}

		if want := "virtual machine with name gitlab-job-2 not found"; err.Error() != want {
			t.Fatalf("expected error %q, got %q", want, err.Error())
		}
	})
}

func TestGetByName_ListFailureIsNotErrNotFound(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	})

	_, _, err := client.VirtualMachines().GetByName(context.Background(), "project", "gitlab-job-1")

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a 401 APIError, got %v", err)
	}

	if errors.Is(err, ErrNotFound) {
		t.Fatal("a failed list must not look like a missing virtual machine")
	}
}
