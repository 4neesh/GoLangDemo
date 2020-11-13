package person

import (
	"advanced-web-service/dataStore"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var personList []Person

func getPersonList() ([]Person, error) {
	results, err := dataStore.DbConn.Query(`SELECT Id, Name, Age FROM people`)
	if err != nil {

		log.Fatal(err)
	}

	defer results.Close()
	personList := make([]Person, 0)
	for results.Next() {
		var person Person
		results.Scan(&person.Id, &person.Name, &person.Age)
		personList = append(personList, person)
	}
	return personList, nil
}

func PersonHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		personList, err := getPersonList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		personJson, err := json.Marshal(personList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(personJson)
	case http.MethodPost:
		var newPerson Person

		personDetails, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(personDetails, &newPerson)
		if newPerson.Id != 0 {
			w.WriteHeader((http.StatusBadRequest))
			return
		}
		newPerson.Id = getNextId()

		personList = append(personList, newPerson)
		w.WriteHeader(http.StatusCreated)
		return
	}

}

func getNextId() int {
	highestId := -1
	for _, product := range personList {
		if highestId < product.Id {
			highestId = product.Id
		}
	}
	return highestId + 1
}

func SinglePersonHandler(w http.ResponseWriter, r *http.Request) {

	urlPathSegments := strings.Split(r.URL.Path, "people/")
	personId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	person, listItemIndex := findProductById(personId)
	if person == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		personJson, err := json.Marshal(person)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(personJson)
	case http.MethodDelete:
		personList = remove(personList, listItemIndex)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodPut:
		var updatedPerson Person
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedPerson)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedPerson.Id != personId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		person = &updatedPerson
		personList[listItemIndex] = *person
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func remove(s []Person, i int) []Person {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findProductById(personId int) (*Person, int) {
	personList, err := getPersonList()
	if err != nil {
		log.Fatal(err)
	}
	for i, person := range personList {
		if person.Id == personId {
			return &person, i
		}
	}
	return nil, 0
}
