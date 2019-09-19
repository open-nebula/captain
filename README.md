[![Build Status](https://travis-ci.org/open-nebula/captain.svg?branch=master)](https://travis-ci.org/open-nebula/captain)

# captain
Manage the docker containers on a single machine.

## Implementation
The captain initiates a socket connection to a spinner.

### Current Assumption for Scope
* The captain knows a spinner to connect to
* The only interactions are to start or kill
* All images are public (as in dockerhub)

### Changes in the Work
* Captain find spinner by beacon
* More complex container interactions
  * Resource limits
  * Container monitoring
  * Storage
  * Container port connections
* Images housed within the system itself
