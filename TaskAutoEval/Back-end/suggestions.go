package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "os"
)

type FormData struct {
    Task1     string `json:"task1"`
    Task2     string `json:"task2"`
    Tools1    string `json:"tools1"`
    Tracking1 string `json:"tracking1"`
    Pain1     string `json:"pain1"`
    Pain2     string `json:"pain2"`
    Goals1    string `json:"goals1"`
}

type ApiResponse struct {
    Suggestions []string `json:"suggestions"`
}

func main() {
    http.HandleFunc("/", formHandler)
    http.HandleFunc("/submit", submitHandler)

    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    // Serve the HTML form
    if r.Method == http.MethodGet {
        tmpl := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <link rel="stylesheet" href="TAEForm.css">
            <title>Task Automation Evaluation Form</title>
        </head>
        <body>
            <div class="form-container">
                <h1>Task Automation Evaluation Form</h1>
                <form id="evaluation-form" action="/submit" method="POST">
                    <fieldset>
                        <legend>Task Identification</legend>
                        <label for="task1">What specific tasks or processes are you currently performing that you believe could be automated? Please list them.</label>
                        <textarea id="task1" name="task1" required></textarea>
                    
                        <label for="task2">Which tasks in your workflow consume the most time? Please describe them.</label>
                        <textarea id="task2" name="task2" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Current Tools and Processes</legend>
                        <label for="tools1">What tools or software are you currently using to manage these tasks? Please specify.</label>
                        <textarea id="tools1" name="tools1" required></textarea>
                    
                        <label for="tracking1">How do you currently track the progress of these tasks? Describe your method.</label>
                        <textarea id="tracking1" name="tracking1" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Pain Points</legend>
                        <label for="pain1">What challenges or frustrations do you encounter with your current task management process? Please elaborate.</label>
                        <textarea id="pain1" name="pain1" required></textarea>
                    
                        <label for="pain2">Are there any repetitive tasks that you find particularly tedious or prone to errors? If so, please describe.</label>
                        <textarea id="pain2" name="pain2" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Goals and Outcomes</legend>
                        <label for="goals1">What are your primary goals for automating these tasks? (e.g., saving time, reducing errors, improving efficiency)</label>
                        <textarea id="goals1" name="goals1" required></textarea>
                    </fieldset>

                    <button type="submit">Submit Feedback</button>
                </form>
            </div>
        </body>
        </html>
        `
        t, _ := template.New("form").Parse(tmpl)
        t.Execute(w, nil)
    }
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        // Parse the form data
        r.ParseForm()
        data := FormData{
            Task1:     r.FormValue("task1"),
            Task2:     r.FormValue("task2"),
            Tools1:    r.FormValue("tools1"),
            Tracking1: r.FormValue("tracking1"),
            Pain1:     r.FormValue ("pain1"),
            Pain2:     r.FormValue("pain2"),
            Goals1:    r.FormValue("goals1"),
        }

        // Send the data to an external API
        suggestions, err := getAutomationSuggestions(data)
        if err != nil {
            http.Error(w, "Failed to get suggestions", http.StatusInternalServerError)
            return
        }

        // Display the suggestions
        fmt.Fprintf(w, "Form submitted successfully!\n")
        fmt.Fprintf(w, "Task 1: %s\n", data.Task1)
        fmt.Fprintf(w, "Task 2: %s\n", data.Task2)
        fmt.Fprintf(w, "Tools: %s\n", data.Tools1)
        fmt.Fprintf(w, "Tracking Method: %s\n", data.Tracking1)
        fmt.Fprintf(w, "Pain Points: %s\n", data.Pain1)
        fmt.Fprintf(w, "Repetitive Tasks: %s\n", data.Pain2)
        fmt.Fprintf(w, "Goals: %s\n", data.Goals1)
        fmt.Fprintf(w, "Automation Suggestions: %v\n", suggestions)
    }
}

func getAutomationSuggestions(data FormData) ([]string, error) {
    apiUrl := "https://task-automation-evaluation.vercel.app/api/suggestions" // Corrected API endpoint
    apiToken := os.Getenv("API_TOKEN") // Use environment variable for API token

    query := fmt.Sprintf("Give the best automation suggestions based on the answers: %s, %s, %s, %s, %s, %s, %s",
        data.Task1, data.Task2, data.Tools1, data.Tracking1, data.Pain1, data.Pain2, data.Goals1)

    requestBody, err := json.Marshal(map[string]string{"query": query})
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiToken) // Add the API token to the request header

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
    }

    var apiResponse ApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        return nil, err
    }

    return apiResponse.Suggestions, nil
} 