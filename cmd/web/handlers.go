package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/freischarler/unityapp/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.

// Change the signature of the home handler so it is defined as a method agains
// *application.

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		app.notFound(w) //http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Create an instance of a templateData struct holding the slice of
	// snippets.
	data := &templateData{Snippets: s}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"../../ui/html/home.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	//ts, err := template.ParseFiles("../../ui/html/home.page.tmpl")

	//Notice that we can pass the slice of file as a variadic parameter?

	ts, err := template.ParseFiles(files...)
	if err != nil {

		// Because the home handler function is now a method against applicatio
		// it can access its fields, including the error logger. We'll write the
		// message to this instead of the standard logger.
		//app.errorLog.Println(err.Error()) //log.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)

		app.serverError(w, err)
		return
	}

	// We then use the Execute() method on the template set to write the templa
	// content as the response body. The last parameter to Execute() represents
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, data)
	if err != nil {
		// Also update the code here to use the error logger from the applicatio
		// struct.
		//app.errorLog.Println(err.Error()) //log.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)

		app.serverError(w, err)
	}
}

// Add a showSnippet handler function.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404
	// not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) //http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our respo
	// and write it to the http.ResponseWriter.Let’s try this out.

	//fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// Initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partial that we made earlier.
	files := []string{
		"../../ui/html/show.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter.
	err = ts.Execute(w, data) //err = ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}

	// Write the snippet data as a plain-text HTTP response body.
	//fmt.Fprintf(w, "%v", s)
}

// Add a createSnippet handler function.

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	// If it's not, use the w.WriteHeader() method to send a 405 status code and
	// the w.Write() method to write a "Method Not Allowed" response body. We
	// then return from the function so that the subsequent code is not executed
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))

		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("Create a new snippet..."))
}

/*Sometimes you might want to serve a single file from within a handler.
//For this there’s the http.ServeFile() function

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../ui/static/file.zip")
}
*/
