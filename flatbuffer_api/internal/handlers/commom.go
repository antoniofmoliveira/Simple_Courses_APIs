package handlers

import (
	"net/http"

	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
	flatbuffers "github.com/google/flatbuffers/go"
)

const octetStream = "application/octet-stream"

func sendFlatBufferMessage(w http.ResponseWriter, message string, httpStatus int) {
	fbBuilder := flatbuffers.NewBuilder(0)
	fbMessage := fbBuilder.CreateString(message)
	fb.MessageStart(fbBuilder)
	fb.MessageAddIsSuccess(fbBuilder, false)
	fb.MessageAddMessage(fbBuilder, fbMessage)
	fbMessageOutput := fb.MessageEnd(fbBuilder)
	fbBuilder.Finish(fbMessageOutput)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(httpStatus)
	w.Write(fbBuilder.FinishedBytes())
}
