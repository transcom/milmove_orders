# Orders API Client

The Orders API Client is a way to connect to and interact with the Orders API on the Orders Server

## Dependencies

**orders-api-client** is written in [Go](https://golang.org/). Aside from go, you will need:

- [GNU Make](https://www.gnu.org/software/make/)

Acquiring and installing these is left as an exercise for the reader.

## Building

To build the client run: `make bin/orders-api-client`

## Usage

There are several commands with the Orders API Client. To list them run:

```sh
$ orders-api-client -h
```

The client can connect to the local development server using either the local development certificates or via a
CAC signed by the same certificate authority as the development server's certificate. As an example, a simple command
using the certificates located in `config/tls/devlocal-mtls.{cer,key}` would be this:

```sh
$ orders-api-client get-orders --orders-uuid ${UUID}
```

For a full command in the devlocal environment you would use:

```sh
$ orders-api-client --hostname orderslocal --port 7443 --insecure --certpath ./config/tls/devlocal-faux-all-orders.cer --keypath ./config/tls/devlocal-faux-all-orders.key get-orders --orders-uuid ${UUID}
```

Using a CAC can be done outside of the docker development environment. You would use these options:

```sh
$ orders-api-client --cac --hostname orderslocal --port 7443 --insecure get-orders --orders-uuid ${UUID}
```

Additional options to configure the client can be listed in the command's help text.

In the following subsections the commands will be explained in more detail.

### Get Orders

The `orders-api-client get-orders` command uses a unique identifier or UUID to fetch orders information from the
Orders API. Example usage is as follows:

```sh
$ orders-api-client get-orders --orders-uuid ${UUID}
```

As output the command will provide a record as returned from the server.

### Get Orders Count

The `orders-api-client get-orders-count` command uses an issuer to fetch the number of orders ingested from the
Order API. This can be used as a quick way to validate that records have been uploaded. Example usage is as follows:

```sh
$ orders-api-client get-orders-count --issuer navy
```

The output is a JSON structure with the count and the issuer:

```json
{"count":7,"issuer":"navy"}
```

### Post Revisions

**NOTE**: Currently the client only support posting revisions for **Navy** orders.

The `orders-api-client post-revisions` command ingests CSV files containing orders information and makes updates to the
database via the Orders API. Example usage is as follows:

```sh
$ orders-api-client post-revisions --csv-file orders.csv
```

As output the command will provide a record as returned from the server for each row in the CSV file.

#### Input format

The comma-delimited CSV files should contain **all** of the following columns:

| Column Name | Description |
| ----------- | ----------- |
| SSN (obligation) | If 9 digits, Social Security Number<br>If 10 digits, EDIPI |
| TAC | Household Goods (HHG) Transportation Account Code (TAC) |
| Order Create/Modification Date | Orders date, in Excel date format (Day 1 = Dec 31, 1899) |
| Order Modification Number | Number of modifications made by an Orders Writing System, such as EAIS, OAIS, or NMCMPS. |
| Obligation Modification Number | Number of modifications made manually via POEMS. |
| Obligation Status Code | **D**: Cancel Obligation, effectively rescinding these Orders<br>**N**: Initial Mod - amended Orders<br>**P**: Initial Obligation - new Orders|
| Obligation Multi-leg Code | Indicates whether either endpoint is TDY.<br>**0** - Perm to Perm<br>**1** - Perm to Temp<br>**5** - Temp to Temp<br>**9** - Temp to Perm |
| CIC Purpose Information Code (OBLGTN) | Purpose of the Orders, maps to Orders type |
| Paygrade | Three-character DoD Paygrade, e.g., E05, W02, O10 |
| Rank Classification  Description | The Navy rank or rating |
| Service Member Name | The sailor's name, in the format LASTNAME,FIRSTNAME (optional MI) (optional suffix) |
| Detach UIC | Unit Identification Code (UIC) of the detaching activity |
| Detach UIC Home Port | Home port of the detaching activity |
| Detach UIC City Name | Detaching activity city |
| Detach State Code | Detaching activity state |
| Detach Country Code | Detaching activity country |
| Ultimate Estimated Arrival Date | Report No Later Than Date |
| Ultimate UIC | Unit Identification Code (UIC) of the ultimate activity |
| Ultimate UIC Home Port | Home port of the ultimate activity |
| Ultimate UIC City Name | Ultimate activity city |
| Ultimate State Code | Ultimate activity state |
| Ultimate Country Code | Ultimate activity country |
| Entitlement Indicator | If 'Y', then this is a 'Cost Order' with obligated moving expenses. If 'N', then this is a 'No Cost Order'. |
| Count of Dependents Participating in Move (STATIC) | Number of sailor's dependents; needed to determine the correct weight entitlement |
| Count of Intermediate Stops (STATIC) | Number of intermediate activities. If greater than 0, then this move has TDY en route. |
| Primary SDN | The Commercial Travel (CT) Standard Document Number (SDN), **used as the unique Orders number** |

Columns that do not start with the above headers are ignored.

#### Orders number

On printed Navy Orders, the BUPERS Orders number is originally formatted as "`<Order Control Number> <SSN>`", for example, "`3108 000-12-3456`". It would be unique (because of the SSN), except that it’s possible for a set of orders to be cut on the same day 10 years later for the same sailor, resulting in a collision.

Because the BUPERS Orders Number contains PII (the SSN) and could potentially not be unique (because it only allows a single digit for the year), the client uses the Primary SDN (aka the Commercial Travel SDN) instead. Similarly, Marine Corps orders also use the CT SDN as the unique Orders number.

#### Modification number interpretation

The Orders API has a sequence number to indicate the chronology of amendments to a set of Orders. The input, however, has two modification number fields, which track the modification count from different systems. Fortunately, these two fields increment atomically, and never decrement.

Therefore, the sequence number is simply the sum of these two numbers.

#### Orders type

To determine the effective orders type, lookup the CIC Purpose Information Code and community (enlisted or officer) in the following table.

| N_CIC_PURP | Enlisted / Officer | Description | Effective Orders Type |
| ---------- | ------------------ | ----------- | --------------------- |
| 0 | Officer | IPCOT In-place consecutive overseas travel | ipcot |
| 1 | Enlisted | IPCOT In-place consecutive overseas travel | ipcot |
| 2 | Officer | Accession Travel | accession |
| 3 | Officer | Training Travel | training |
| 4 | Officer | Operational Travel | operational |
| 5 | Officer | Separation Travel | separation |
| 6 | Officer | Organized Unit/Homeport Change | unit-move |
| 7 | Officer | Emergency Non-member Evac | emergency-evac |
| 8 | Enlisted | Overseas Tour Extension Incentive Program (OTEIP) | oteip |
| 9 | Enlisted | NAVCAD (Naval Cadet) Training | training |
| A | Enlisted | Accession Travel Recruits | accession |
| B | Enlisted | Non-recruit Accession Travel | accession |
| C | Enlisted | Training Travel | training |
| D | Enlisted | Operational Travel | operational |
| E | Enlisted | Separation Travel | separation |
| F | Enlisted | Organized Unit/Homeport Change | unit-move |
| G | Enlisted | Midshipman Accession Travel | accession |
| H | Both | Special Purpose Reimbursable | special-purpose |
| I | Enlisted | NAVCAD(Naval Cadet) Accession | accession |
| J | Enlisted | Accession Travel Recruits | accession |
| K | Enlisted | Non-recruit Accession Travel | accession |
| L | Enlisted | Training Travel | training |
| M | Enlisted | Rotational Travel | rotational |
| N | Enlisted | Separation Travel | separation |
| O | Enlisted | Organized Unit/Homeport Change | unit-move |
| P | Enlisted | Midshipman Separation Travel | separation |
| Q | Officer | Misc. Rotational Non-member | rotational |
| R | Enlisted | Misc. Operational Non-member | operational |
| S | Officer | Accession Travel | accession |
| T | Officer | Training Travel | training |
| U | Officer | Rotational Travel | rotational |
| V | Officer | Separation Travel | separation |
| W | Officer | Organized Unit/Homeport Change | unit-move |
| X | Enlisted | EMERGENCY NON-MEMBER EVACS | emergency-evac |
| X | Officer | Misc. Rotational Non-member | rotational |
| Y | Enlisted | Misc. Rotational Non-member | rotational |
| Z | Enlisted | NAVCAD(Naval Cadet) Separation | separation |
