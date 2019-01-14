# Bad Licenses

Detects bad licenses using two processors:

* Processor 1 tracks bad licenses by simply storing messages in "configure-licenses" as state.
* Processor 2 consumes startTrip events and joins (lookup) with the bad-license state, notifying if a strip with a blacklisted license was started.
