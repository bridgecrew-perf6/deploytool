OrdererOrgs:{{range $key,$value:= .ordList}}
  - Name: ord{{$key}}
    Domain: ord{{$key}}.{{$.domain}}
    EnableNodeOUs: true
    Template:
      Count: {{$value}}{{end}}

PeerOrgs:{{range $key,$value:= .orgList}}
  - Name: org{{$key}}
    Domain: org{{$key}}.{{$.domain}}
    EnableNodeOUs: true
    Template:
      Count: {{$value}}
      SANS:
        - localhost
    Users:
      Count: 2{{end}}
