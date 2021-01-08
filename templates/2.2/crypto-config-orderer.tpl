OrdererOrgs:
  - Name: {{.orgId}}
    Domain: {{.orgId}}.{{.domain}}
    EnableNodeOUs: true
    Template:
      Count: {{.nodeCounts}}

