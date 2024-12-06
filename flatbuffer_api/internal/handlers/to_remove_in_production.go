

package handlers

import "net/http"

// TODO! only for test - remove in production
func (h *CategoryHandler) CategoriesError(w http.ResponseWriter, r *http.Request) {
	sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
}
