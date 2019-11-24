# What happens when you type address in browser and hit return

## DNS: convert host name to IP address.
DNS servers usually provided by ISP (internet service provider).
DNS servers connect to ***, which connect to domain server (e.g. domain.com).

DNS server returns IP address.

Browser sets up TCP connection (using IP address)

Browser requests data from IP address.
* Application
* Presentation
* Session
* Transport (TCP)
* Networ (IP)
* Data (Ethernet)
* Physical

Browser return status code and (possible) info
* 100's: informations
* 200's: normal (returns html packet from server)
* 300's: redirect (particularly code:303)
* 400's: client error (e.g. page not found)
* 500's: server error
Use 303 to redirect to secure HTTPS connection

Browser renders HTML
