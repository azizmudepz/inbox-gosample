package message

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drone/routes"
	"github.com/julienschmidt/httprouter"
	"github.com/rtulus/inbox-gosample/src/conf"
)

type InboxData struct {
	InboxID    int64
	UserID     int64
	Status     int
	ReadStatus int
}

func ReadInbox(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"success": 0,
	}

	db := conf.DB.Database
	userID := ps.ByName("user_id")
	// userIDInt, _ := strconv.ParseInt(userID, 10, 64)

	query := fmt.Sprintf(`
	    SELECT inbox_id, user_id, status, read_status
        FROM ws_inbox_message
	    WHERE user_id = %d`,
		userID)
	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		routes.ServeJson(w, response)
		return
	}
	defer rows.Close()

	inboxList := []InboxData{}
	for rows.Next() {
		inbox := InboxData{}

		errScan := rows.Scan(
			&inbox.InboxID,
			&inbox.UserID,
			&inbox.Status,
			&inbox.ReadStatus,
		)
		if errScan != nil {
			routes.ServeJson(w, response)
			return
		}
		inboxList = append(inboxList, inbox)
	}

	response = map[string]interface{}{
		"success": 1,
		"data":    inboxList,
	}
	routes.ServeJson(w, response)
	return
}
