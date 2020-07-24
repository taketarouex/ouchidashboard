# ouchi-dashboard

display and store data related to `ouchi` life.

## diagram

@startuml
!include ./common.puml
!include ./GCP Icons/Products & Services/Storage & Databases/Cloud Bigtable.puml
!include ./GCP Icons/Products & Services/Big Data/BigQuery.puml

GCP_CloudBigtable(foo, "My BigTable")
GCP_BigQuery(bar, "Data Processing")

foo -> bar

@enduml
