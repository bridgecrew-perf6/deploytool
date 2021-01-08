Organizations:{{range $key,$value:= .ordList}}
    - &OrdererOrg{{$key}}
        Name: {{$key}}
        ID: {{$key}}
        MSPDir: crypto-config/ordererOrganizations/{{$key}}.{{$.domain}}/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('{{$key}}.member')"
            Writers:
                Type: Signature
                Rule: "OR('{{$key}}.member')"
            Admins:
                Type: Signature
                Rule: "OR('{{$key}}.admin')"{{end}}
    {{range $key,$value:= .orgList}}
    - &Org{{$key}}
        Name: {{$key}}
        ID: {{$key}}
        MSPDir: crypto-config/peerOrganizations/{{$key}}.{{$.domain}}/msp
        Policies:
             Readers:
                 Type: Signature
                 Rule: "OR('{{$key}}.admin', '{{$key}}.peer', '{{$key}}.client')"
             Writers:
                 Type: Signature
                 Rule: "OR('{{$key}}.admin', '{{$key}}.client')"
             Admins:
                 Type: Signature
                 Rule: "OR('{{$key}}.admin')"
             Endorsement:
                 Type: Signature
                 Rule: "OR('{{$key}}.peer')"
        AnchorPeers:{{range $index,$peer:= $.peers}} {{if eq $peer.orgId $key}} {{if eq $peer.id "0"}}
            - Host: peer0.{{$peer.orgId}}.{{$.domain}}
              Port: {{$peer.externalPort}}{{end}}{{end}}{{end}}{{end}}

Capabilities:
    Channel: &ChannelCapabilities
        V2_0: true
    Orderer: &OrdererCapabilities
        V2_0: true
    Application: &ApplicationCapabilities
        V2_0: true

Application: &ApplicationDefaults
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: "MAJORITY Endorsement"
        Endorsement:
            Type: ImplicitMeta
            Rule: "ANY Endorsement"

    Capabilities:
        <<: *ApplicationCapabilities
Orderer: &OrdererDefaults
    OrdererType: etcdraft
    BatchTimeout: {{.batchTime}}
    BatchSize:
        MaxMessageCount: {{.batchSize}}
        AbsoluteMaxBytes: 98 MB
        PreferredMaxBytes: {{.batchPreferred}}
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        BlockValidation:
            Type: ImplicitMeta
            Rule: "ANY Writers"

Channel: &ChannelDefaults
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
    Capabilities:
        <<: *ChannelCapabilities

Profiles:
    OrgsChannel:
        Consortium: SampleConsortium
        <<: *ChannelDefaults
        Application:
            <<: *ApplicationDefaults
            Organizations:{{range $key,$value:= .orgList}}
                - *Org{{$key}}{{end}}
            Capabilities:
                <<: *ApplicationCapabilities

    SampleMultiNodeEtcdRaft:
        <<: *ChannelDefaults
        Capabilities:
            <<: *ChannelCapabilities
        Orderer:
            <<: *OrdererDefaults
            OrdererType: etcdraft
            EtcdRaft:
                Consenters:{{range $index,$orderer:= .orderers}}
                - Host: orderer{{$orderer.id}}.{{$orderer.orgId}}.{{$.domain}}
                  Port: {{$orderer.externalPort}}
                  ClientTLSCert: crypto-config/ordererOrganizations/{{$orderer.orgId}}.{{$.domain}}/orderers/orderer{{$orderer.id}}.{{$orderer.orgId}}.{{$.domain}}/tls/server.crt
                  ServerTLSCert: crypto-config/ordererOrganizations/{{$orderer.orgId}}.{{$.domain}}/orderers/orderer{{$orderer.id}}.{{$orderer.orgId}}.{{$.domain}}/tls/server.crt{{end}}
            Addresses:{{range $index,$orderer:= .orderers}}
                - orderer{{$orderer.id}}.{{$orderer.orgId}}.{{$.domain}}:{{$orderer.externalPort}}{{end}}
            Organizations:{{range $key,$value:= .ordList}}
            - *OrdererOrg{{$key}}{{end}}
            Capabilities:
                <<: *OrdererCapabilities
        Application:
            <<: *ApplicationDefaults
            Organizations:{{range $key,$value:= .ordList}}
            - <<: *OrdererOrg{{$key}}{{end}}
        Consortiums:
            SampleConsortium:
                Organizations:{{range $key,$value:= .orgList}}
                - *Org{{$key}}{{end}}

