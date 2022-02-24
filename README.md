# Unisearcher API

An API that provides a more detailed list of universities around the world


### Using The Hosted API
This is by far the easiest way of using the API.  
There are three endpoins available, where two of them are used to retrieve information and the last one is used for diagnostics.   

* https://unisearcher.herokuapp.com/unisearcher/v1/uniinfo/
* https://unisearcher.herokuapp.com/unisearcher/v1/neighbourunis/
* https://unisearcher.herokuapp.com/unisearcher/v1/diag/


#### Uniinfo
The first endpoint is used to search for universities by name, either by full name or by partial name.  

**Request**
```
METHOD: GET
Path: https://unisearcher.herokuapp.com/unisearcher/v1/uniinfo/{:partial_or_complete_university_name}/
```

**Examples:**

https://unisearcher.herokuapp.com/unisearcher/v1/uniinfo/molde%20university%20college  
https://unisearcher.herokuapp.com/unisearcher/v1/uniinfo/norwegian

**Response**
* Content type: application/json

```
[
  {
    "name": "Molde University College",
    "country": "Norway",
    "isocode": "NO",
    "webpages": [
      "http://www.himolde.no/"
    ],
    "languages": {
      "nno": "Norwegian Nynorsk",
      "nob": "Norwegian Bokmål",
      "smi": "Sami"
    },
    "map": "https://www.openstreetmap.org/relation/2978650"
  }
]
```

#### Neighbourunis
The second endpoint is used to search for universities in neighbouring countries sharing a partial name. 

**Request**
```
METHOD: GET
Path: https://unisearcher.herokuapp.com/unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}
```
{:country_name} refers to the English name for the country that is the basis (basis country) of the search of unis with the same name in neighbouring countries.  

{:partial_or_complete_university_name} is the partial or complete university name, for which universities with similar name are sought in neighbouring countries  

{?limit={:number}} is an optional parameter that limits the number of universities in bordering countries (number) that are reported.  

**Examples:**  
https://unisearcher.herokuapp.com/unisearcher/v1/neighbourunis/norway/science   
https://unisearcher.herokuapp.com/unisearcher/v1/neighbourunis/norway/science?limit=3  

**Response**
* Content type: application/json

```
[
  {
    "name": "Central Ostrobothnia University of Applied Sciences",
    "country": "Finland",
    "isocode": "FI",
    "webpages": [
      "http://www.cou.fi/"
    ],
    "languages": {
      "fin": "Finnish",
      "swe": "Swedish"
    },
    "map": "openstreetmap.org/relation/54224"
  },
  {
    "name": "Diaconia University of Applied Sciences",
    "country": "Finland",
    "isocode": "FI",
    "webpages": [
      "http://www.diak.fi/"
    ],
    "languages": {
      "fin": "Finnish",
      "swe": "Swedish"
    },
    "map": "openstreetmap.org/relation/54224"
  },
  {
    "name": "EVTEK University of Applied Sciences",
    "country": "Finland",
    "isocode": "FI",
    "webpages": [
      "http://www.evtek.fi/"
    ],
    "languages": {
      "fin": "Finnish",
      "swe": "Swedish"
    },
    "map": "openstreetmap.org/relation/54224"
  }
]
```

#### Diag
The last endpoint is used to get diagnostics of the unisearcher API  

**Request**
```
METHOD: GET
Path: https://unisearcher.herokuapp.com/unisearcher/v1/diag/
```

**Response**
* Content type: application/json  

```
{
  "universityAPI": "200 OK",
  "countryAPI": "200 OK",
  "version": "v1",
  "uptime": "1393.706649168"
}
```

