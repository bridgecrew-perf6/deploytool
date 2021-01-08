PeerOrgs:
  - Name:  {{.orgId}}
    Domain: {{.orgId}}.{{.domain}}
    EnableNodeOUs: true
    Template:
      Count: {{.nodeCounts}}
    Users:
      Count: 2
