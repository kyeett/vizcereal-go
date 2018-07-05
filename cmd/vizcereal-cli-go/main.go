package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func GetJson(c int) string {
	return `{
  "renderer": "global",
  "name": "edge",
  "nodes": [
    {
      "renderer": "region",
      "name": "INTERNET",
      "class": "normal"
    },
    {
      "renderer": "region",
      "name": "us-west-1",
      "maxVolume": 50000,
      "class": "normal",
      "nodes": []
    },
    {
      "renderer": "region",
      "name": "us-east-1",
      "maxVolume": 50000,
      "class": "normal",
      "updated": 1466838546805,
      "nodes": [
        {
          "name": "INTERNET",
          "renderer": "focusedChild",
          "class": "normal"
        },
        {
          "name": "proxy-prod",
          "renderer": "focusedChild",
          "class": "normal"
        },
        {
          "name": "end",
          "renderer": "focusedChild",
          "class": "normal"
        },
        {
          "name": "proxy-prod2",
          "renderer": "focusedChild",
          "class": "normal",
          "notices": [
            {
              "title": "Notice about something",
              "link": "http://link/to/relevant/thing",
              "severity": 1
            },{
              "title": "Notice about something2",
              "link": "http://link/to/relevant/thing2",
              "severity": 2
            }
          ]
        }
      ],
      "connections": [
        {
          "source": "end",
          "target": "INTERNET",
          "metrics": {
            "danger": 116.524,
            "normal": 598.906
          },
          "class": "normal"
        },{
          "source": "INTERNET",
          "target": "proxy-prod",
          "metrics": {
            "danger": 116.524,
            "normal": 598.906
          },
          "class": "normal"
        },{
          "source": "INTERNET",
          "target": "proxy-prod2",
          "metrics": {
            "danger": 116.524,
            "normal": 198.906
          },
          "class": "normal"
        },{
          "source": "proxy-prod2",
          "target": "end",
          "metrics": {
            "danger": 116.524,
            "normal": 158.906
          },
          "class": "normal"
        },{
          "source": "proxy-prod",
          "target": "proxy-prod2",
          "metrics": {
            "danger": 116.524,
            "normal": 558.906
          },
          "class": "normal"
        }
      ]
    }
  ],
  "connections": [
    {
      "source": "INTERNET",
      "target": "us-east-1",
      "metrics": {
        "normal": ` + strconv.Itoa(c) + `,
        "danger": 92.37
      },
      "notices": [
      ],
      "class": "normal"
    },
    {
      "source": "INTERNET",
      "target": "us-west-1",
      "metrics": {
        "danger": ` + strconv.Itoa(100000-c) + `,
        "normal": 92.37
      },
      "notices": [
      ],
      "class": "normal"
    }
  ]
}`
}

var i int = 1

func editHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, struct {
		Json string
	}{
		Json: GetJson(10),
	})
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	i = ((i) * 2)
	if i > 200000 {
		i = 1
		fmt.Println("Too big")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, GetJson(i))
}

func main() {
	fmt.Println("test2")
	fs := http.FileServer(http.Dir(""))
	http.Handle("/static/", fs)
	http.HandleFunc("/", editHandler)
	http.HandleFunc("/test/", serveHandler)
	http.ListenAndServe(":8081", nil)
}
