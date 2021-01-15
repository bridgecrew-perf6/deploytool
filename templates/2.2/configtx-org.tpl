Organizations:
    - &{{.msp_id}}
        Name: {{.msp_id}}
        ID: {{.msp_id}}
        MSPDir: {{.msp_path}}
        Policies:
             Readers:
                 Type: Signature
                 Rule: "OR('{{.msp_id}}.admin', '{{.msp_id}}.peer', '{{.msp_id}}.client')"
             Writers:
                 Type: Signature
                 Rule: "OR('{{.msp_id}}.admin', '{{.msp_id}}.client')"
             Admins:
                 Type: Signature
                 Rule: "OR('{{.msp_id}}.admin')"
             Endorsement:
                 Type: Signature
                 Rule: "OR('{{.msp_id}}.peer')"
        AnchorPeers:
            - Host: {{.peer_address}}
              Port: {{.peer_port}}




