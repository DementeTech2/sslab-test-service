define({ "api": [
  {
    "type": "get",
    "url": "/api/analyze/:domain?[dont_wait]",
    "title": "Start to analyze a domain",
    "name": "AnalyzeDomain",
    "group": "Domains",
    "version": "1.0.0",
    "description": "<p>This is the Description.</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "optional": false,
            "field": "domain",
            "description": "<p>Include the servers of every domain in the response</p>"
          },
          {
            "group": "Parameter",
            "optional": true,
            "field": "dont_wait",
            "description": "<p>Start the analysis async. By default is going to timeout at 60 seconds.</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Initial regular response (dont_wait):",
          "content": "HTTP/1.1 200 OK\n{\n    \"id\"                 : 123456789\n    \"domain\"             : \"google.com\"\n    \"start_time\"         : \"2019-06-17T22:00:21-0000\"\n    \"end_time\"           : \"1970-01-01T00:00:00-0000\"\n    \"status\"             : \"pending\"\n    \"logo\"               : \"\"\n    \"title\"              : \"\"\n    \"ssl_grade\"          : \"\"\n    \"previous_ssl_grade\" : \"\"\n    \"server_changed\"     : false\n    \"is_down\"            : false\n    \"servers\"            : []\n}",
          "type": "json"
        },
        {
          "title": "Regular response without timeout:",
          "content": "HTTP/1.1 200 OK\n{\n    \"id\"                 : 123456789\n    \"domain\"             : \"google.com\"\n    \"start_time\"         : \"2019-06-17T22:00:21-0000\"\n    \"end_time\"           : \"2019-06-17T22:15:21-0000\"\n    \"status\"             : \"ready\"\n    \"logo\"               : \"https://google.com/logog.png\"\n    \"title\"              : \"Google com\"\n    \"ssl_grade\"          : \"a\"\n    \"previous_ssl_grade\" : \"\"\n    \"server_changed\"     : false\n    \"is_down\"            : false\n    \"servers\"            : []\n}",
          "type": "json"
        },
        {
          "title": "Timeout:",
          "content": "HTTP/1.1 408 Request Timeout\nStill running, call it later",
          "type": "json"
        }
      ]
    },
    "filename": "src/ssllabtestservice/server/server.go",
    "groupTitle": "Domains"
  },
  {
    "type": "get",
    "url": "/api/fetch-web-data/:domain",
    "title": "Fetch Domain Web Data",
    "name": "FetchDomainData",
    "group": "Domains",
    "version": "1.0.0",
    "description": "<p>This is the Description.</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "optional": false,
            "field": "domain",
            "description": "<p>The domain to analyze</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Good domain:",
          "content": "HTTP/1.1 200 OK\n{\n  \"Domain\": \"fake.foo\",\n  \"Title\": \"The fake page\",\n  \"Logo\": \"https://fake.foo/my-logo.png\",\n  \"IsDown\": false\n}",
          "type": "json"
        },
        {
          "title": "Bad domain:",
          "content": "HTTP/1.1 200 OK\n{\n  \"Domain\": \"fake.down.foo\",\n  \"Title\": \"NOT_FOUND\",\n  \"Logo\": \"NOT_FOUND\",\n  \"IsDown\": true\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "BadRequest",
            "description": "<p>The <code>domain</code> is not valid</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Domain not valid:",
          "content": "HTTP/1.1 400 BadRequest\nInvalid Domain",
          "type": "json"
        }
      ]
    },
    "filename": "src/ssllabtestservice/server/server.go",
    "groupTitle": "Domains"
  },
  {
    "type": "get",
    "url": "/api/domains?[include_servers]",
    "title": "Get all domains analyzed",
    "name": "GetDomains",
    "group": "Domains",
    "version": "1.0.0",
    "description": "<p>This is the Description.</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "optional": true,
            "field": "include_servers",
            "description": "<p>Include the servers of every domain in the response</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": ":",
          "content": "HTTP/1.1 200 OK\n[\n    {\n        \"id\"                 : 123456789,\n        \"domain\"             : \"google.com\",\n        \"start_time\"         : \"2019-06-17T22:00:21-0000\",\n        \"end_time\"           : \"2019-06-17T22:15:21-0000\",\n        \"status\"             : \"ready\",\n        \"logo\"               : \"https://google.com/logog.png\",\n        \"title\"              : \"Google com\",\n        \"ssl_grade\"          : \"a\",\n        \"previous_ssl_grade\" : \"\",\n        \"server_changed\"     : false,\n        \"is_down\"            : false,\n        \"servers\"            : [\n             {\n                 \"id\"          : 9876,\n                 \"revision_id\" : 123456789,\n                 \"ip\"          : \"127.0.0.1\",\n                 \"ssl_grade\"   : \"a\",\n                 \"progress\"    : 100,\n                 \"country\"     : \"us\",\n                 \"owner\"       : \"GoogleInc\"\n             },\n             {\n                 \"id\"          : 9877,\n                 \"revision_id\" : 123456789,\n                 \"ip\"          : \"127.0.0.1\",\n                 \"ssl_grade\"   : \"a\",\n                 \"progress\"    : 100,\n                 \"country\"     : \"us\",\n                 \"owner\"       : \"GoogleInc\"\n             },\n             ....\n        ]\n    },\n    {\n        \"id\"                 : 123456790,\n        \"domain\"             : \"fake.com\",\n        \"start_time\"         : \"2019-06-17T22:00:21-0000\",\n        \"end_time\"           : \"2019-06-17T22:15:21-0000\",\n        \"status\"             : \"error\",\n        \"logo\"               : \"NOT_FOUND\",\n        \"title\"              : \"NOT_FOUND\",\n        \"ssl_grade\"          : \"\",\n        \"previous_ssl_grade\" : \"\",\n        \"server_changed\"     : false,\n        \"is_down\"            : true,\n        \"servers\"            : [],\n    },\n    ....\n]",
          "type": "json"
        }
      ]
    },
    "filename": "src/ssllabtestservice/server/server.go",
    "groupTitle": "Domains"
  }
] });
