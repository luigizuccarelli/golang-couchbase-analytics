let bDebug = false;
let bTrace = false;

// Analytics data schema
var schema = {
  "from": {
    "url":"",
    "pagename":"",
    "pagetype":""
  },
  "to":  {
    "url":"",
    "pagename":"",
    "pagetype":""
  },
  "location": {
    "ip": "",
    "carrier":"",
    "country": {
      "name":"",
      "code":"",
      "capital":""
    }
  },
  "currency": {
    "name":"",
    "code":""
  },
  "event": {
    "type":"",
    "timeonpage":""
  },
  "utm_campaign":"",
  "utm_affiliate":"",
  "utm_medium":"",
  "utm_source":"",
  "utm_content":"",
  "timestamp":"",
  "platform": {
    "appCodeName":"",
    "appName":"",
    "appVersion":"",
    "language":"",
    "os":"",
    "product":"",
    "productSub":"",
    "userAgent":"",
    "vendor":""
  },
  "trackingid":"",
  "creative": {
    "name": "",
    "status": ""
  },
  "effort" : {
    "acquisition_method": "",
    "advantage_campaign_code":"",
    "advantage_description": "",
    "advertisement_name": "",
    "campaign": "",
    "domain" : "",
    "date": "",
    "effort_destination":"",
    "id": "",
    "type":"",
    "journey":"",
    "promocode":"",
    "what_are_you_promoting":"",
    "where_is_the_marketing_going":""
  },
  "journey": {
    "creative_sequence":"",
    "name":"",
    "status":""
  }
};

let data_stream_url = "https://message-producer-trackmate-poc.apps.balt1.okd.14west.io/api/v1/streamdata"

/*
 * function getInfo - does an ip lookup and gets all the relevant meta data
 *
 * @params void
 * @returns void
 */
function getInfo(_callback) {
  xhttp = new XMLHttpRequest();
  xhttp.open("GET", "https://api.ipregistry.co/?key=g4g3ykru1blcd6" , true);
  xhttp.send();

  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      let json = JSON.parse(this.responseText);
      if (bTrace) {
        console.log(this.responseText);
      }
      schema.location.ip = json.ip;
      schema.location.carrier = json.carrier.name; 
      schema.location.country.name = json.location.country.name; 
      schema.location.country.code = json.location.country.code; 
      schema.location.country.capital = json.location.country.capital; 
      schema.currency.name = json.currency.name; 
      schema.currency.code = json.currency.code; 
      _callback();
    } 
  }
}

/*
 * postAnalyticsData to endpoint
 * 
 * @returns - void
 * @params - json payload
 *
 */
function postAnalyticsData(json) {
  let xhttp = new XMLHttpRequest();
  xhttp.open("POST", data_stream_url , true);
  xhttp.send(json);

  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      console.log(JSON.parse(this.responseText));
      return;
    }
  }
}


/*
 * function getJsonFromUrl - utility function that parses the url and creates a key value pair for each parameter
 *
 * @params url (location.href)
 * @returns void
 */
function getJsonFromUrl(url) {
  if(!url) url = location.href;
  var question = url.indexOf("?");
  var hash = url.indexOf("#");
  if(hash==-1 && question==-1) return {};
  if(hash==-1) hash = url.length;
  var query = question==-1 || hash==question+1 ? url.substring(hash) : 
  url.substring(question+1,hash);
  var result = {};
  query.split("&").forEach(function(part) {
    if(!part) return;
    part = part.split("+").join(" "); // replace every + with space, regexp-free version
    var eq = part.indexOf("=");
    var key = eq>-1 ? part.substr(0,eq) : part;
    var val = eq>-1 ? decodeURIComponent(part.substr(eq+1)) : "";
    var from = key.indexOf("[");
    if(from==-1) result[decodeURIComponent(key)] = val;
    else {
      var to = key.indexOf("]",from);
      var index = decodeURIComponent(key.substring(from+1,to));
      key = decodeURIComponent(key.substring(0,from));
      if(!result[key]) result[key] = [];
      if(!index) result[key].push(val);
      else result[key][index] = val;
    }
  });
  return result;
}

function injectParams(url) {
  try {
    let newUrl = "";
    let prefix = "";

    if (url.indexOf("?") > 1) {
      prefix = "&";
    } else {
      prefix = "?";
    }
    pagetype = document.getElementById("pagetype").innerHTML;
    if (pagetype === 'origin') {
      newUrl = url + prefix + "trackingid=" + schema.trackingid;
    } else {
      newUrl = url + prefix + "trackingid=" + schema.trackingid 
             + "&utm_affiliate=" + schema.utm_affiliate
             + "&utm_campaign=" + schema.utm_campaign
             + "&utm_medium=" + schema.utm_medium
             + "&utm_source=" + schema.utm_source
             + "&utm_content=" + schema.utm_content
             + "&pagename=" + schema.to.pagename
             + "&pagetype=" + schema.to.pagetype;
    }
    location.href = newUrl;
  } catch(e) {
    console.error(e);
    location.href = url;
  }
} 

let startDate = new Date();
let elapsedTime = 0;

// Create eventListeners
window.addEventListener('load', function () {
  // wait for ip call to complete
  getInfo(function() { 
    let meta = window.navigator;
    urlParams = getJsonFromUrl(location.href)
    // disable ip and geo location for now
    if (urlParams["custom_referrer"]) {
      schema.from.url = urlParams["custom_referrer"];
    } else {
      schema.from.url = document.referrer;
    }
    schema.to.url = location.href;
    schema.event.type = "load";
    schema.event.timeonpage = 0;
   
    try {
      schema.to.pagename = document.getElementById("pagename").innerHTML; 
      schema.to.pagetype = document.getElementById("pagetype").innerHTML;
      schema.from.pagename = urlParams["pagename"]; 
      schema.from.pagetype = urlParams["pagetype"];
      schema.utm_campaign = urlParams["utm_campaign"];
      schema.utm_affiliate =  urlParams["utm_affiliate"];
      schema.utm_medium = urlParams["utm_medium"];
      schema.utm_source = urlParams["utm_source"];
      schema.utm_content = urlParams["utm_content"];
    } catch(e) {
      console.error(e);
    }
    schema.platform.appCodeName = meta.appCodeName;  
    schema.platform.appName = meta.appName;  
    schema.platform.appVersion = meta.appVersion;  
    schema.platform.language = meta.language;  
    schema.platform.os = meta.platform;  
    schema.platform.product = meta.product;  
    schema.platform.productSub = meta.productSub;  
    schema.platform.userAgent = meta.userAgent;  
    schema.platform.vendor = meta.vendor;  
    schema.timestamp = new Date().getTime();
    if (schema.to.pagetype === 'origin') {
      schema.trackingid = uuidv4();
    } else {
      schema.trackingid = urlParams["trackingid"];
    }
    if(typeof irisPlusData !== "undefined") {
      schema.creative.name = irisPlusData["CREATIVE.name"];
      schema.creative.status = irisPlusData["CREATIVE.status"];
      schema.effort.acquisition_method = irisPlusData["EFFORT.acquisition_method"];
      schema.effort.advantage_campaign_code = irisPlusData["EFFORT.advantage_campaign_code"];
      schema.effort.advantage_description = irisPlusData["EFFORT.advantage_description"];
      schema.effort.advertisement_name = irisPlusData["EFFORT.advertisement_name"];
      schema.effort.campaign = irisPlusData["EFFORT.campaign"];
      schema.effort.domain = irisPlusData["EFFORT.domain"];
      schema.effort.date = irisPlusData["EFFORT.date"];
      schema.effort.effort_destination = irisPlusData["EFFORT.effort_destination"];
      schema.effort.id = irisPlusData["EFFORT.id"];
      schema.effort.type = irisPlusData["EFFORT.type"];
      schema.effort.journey = irisPlusData["EFFORT.journey"];
      schema.effort.promocode = irisPlusData["EFFORT.promocode"];
      schema.effort.what_are_you_promoting = irisPlusData["EFFORT.what_are_you_promoting"];
      schema.effort.where_is_the_marketing_going = irisPlusData["EFFORT.where_is_the_marketing_going"];
      schema.journey.creative_sequence = irisPlusData["JOURNEY.creative_sequence"];
      schema.journey.name = irisPlusData["JOURNEY.name"];
      schema.journey.status = irisPlusData["JOURNEY.status"];
    }

    if (bTrace) {
      JSONstringify(schema);
    }
    pagetype = document.getElementById("pagetype").innerHTML;
    if (!bDebug && pagetype !== 'origin') {
      // post to our data stream
      postAnalyticsData(JSON.stringify(schema));
    }

  });

});

window.addEventListener("beforeunload", function (e) {
  endDate = new Date();
  schema.event.type = "exit";
  schema.event.timeonpage = endDate.getTime() - startDate.getTime();
  schema.timestamp = new Date().getTime();
 
  pagetype = document.getElementById("pagetype").innerHTML;
  if (!bDebug && pagetype !== 'origin') {
    postAnalyticsData(JSON.stringify(schema));
  }

});

// Cookie management
function setCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")  + expires + "; SameSite=None; Secure" + "; path=/";
}

// Get cookies
function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

// Erase cookies
function eraseCookie(name) {   
    document.cookie = name+'=; Max-Age=-99999999;';  
}

// Create unique id 
function uuidv4() {
  return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

/*
 * Utility funcion
 */
function JSONstringify(json) {
    if (typeof json != 'string') {
        json = JSON.stringify(json, undefined, '\t');
    }

    var 
        arr = [],
        _string = 'color:green',
        _number = 'color:darkorange',
        _boolean = 'color:blue',
        _null = 'color:magenta',
        _key = 'color:red';

    json = json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var style = _number;
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                style = _key;
            } else {
                style = _string;
            }
        } else if (/true|false/.test(match)) {
            style = _boolean;
        } else if (/null/.test(match)) {
            style = _null;
        }
        arr.push(style);
        arr.push('');
        return '%c' + match + '%c';
    });
    arr.unshift(json);
    console.log.apply(console, arr);
}

// use https://www.minifier.org/
