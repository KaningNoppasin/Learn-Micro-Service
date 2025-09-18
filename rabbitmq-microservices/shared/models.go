package shared

import "time"

type Message struct {
    ID        string                 `json:"id"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    Source    string                 `json:"source"`
    Step      int                    `json:"step"`
}

type ProcessingResult struct {
    Success   bool                   `json:"success"`
    Message   string                 `json:"message"`
    Data      map[string]interface{} `json:"data"`
    ProcessedBy string               `json:"processed_by"`
}
