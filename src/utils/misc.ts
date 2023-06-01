const crypt = require('crypto');
import https from 'node:https';

export const randomString = (length: number, chars: String) => {
  if (!chars) {
    throw new Error('Argument \'chars\' is undefined');
  }

  const charsLength = chars.length;
  if (charsLength > 256) {
    throw new Error('Argument \'chars\' should not have more than 256 characters'
      + ', otherwise unpredictability will be broken');
  }

  const randomBytes = crypt.randomBytes(length);
  let result = new Array(length);

  let cursor = 0;
  for (let i = 0; i < length; i++) {
    cursor += randomBytes[i];
    result[i] = chars[cursor % charsLength];
  }

  return result.join('');
}

export const randomAsciiString = (length: number) => {
  return randomString(length,
    'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789');
}

export const simpleStringify = (object: any) => {
  // stringify an object, avoiding circular structures
  // https://stackoverflow.com/a/31557814
  let simpleObject: any = {};
  for (let prop in Object.keys(object)){
    if (!object.hasOwnProperty(prop)){
        continue;
    }
    if (prop == 'edupage') {
        continue;
    }
    if (typeof object[prop] == 'function'){
        continue;
    }
    simpleObject[prop] = object[prop];
  }
  return JSON.stringify(simpleObject); // returns cleaned up JSON
};

export const replaceCircular = (obj: any, level = 0, already = new WeakSet()): object | String => {
  switch (typeof obj) {
    case 'object':
      if (!obj)
        return obj
      if (already.has(obj)) {
        return "CIRCULAR"
      }
      already.add(obj)
      if (Array.isArray(obj)) {
        return obj.map(item => replaceCircular(item, level + 1, already))
      }
      const newObj: any = {}
      Object.keys(obj).forEach(key => {
        const val = replaceCircular(obj[key], level + 1, already)
        newObj[key] = val
      })
      already.delete(obj)
      return newObj
    default:
      return obj;
  }
}

export const replaceCircularFast = (val: any, cache: WeakSet<object>): object | String => {

  cache = cache || new WeakSet();

  if (val && typeof(val) === 'object') {
    if (cache.has(val)) return '[Circular]';

    cache.add(val);

    var obj: any = (Array.isArray(val) ? [] : {});
    for(var idx in val) {
      obj[idx] = replaceCircular(val[idx], 0, cache);
    }

    cache.delete(val);
    return obj;
  }

  return val;
};

export const sendNotification = (data: object) => {
  let headers = {
    "Content-Type": "application/json; charset=utf-8",
    "Authorization": "Basic ZTE3ZDE4YjctZWExYy00NTJkLTkwODgtYTAwNDYzMmM5NzAz"
  };
  
  let options = {
    host: "onesignal.com",
    port: 443,
    path: "/api/v1/notifications",
    method: "POST",
    headers: headers
  };

  let req = https.request(options, function(res) {  
    res.on('data', function(data) {
      console.log("Response:");
      console.log(JSON.parse(data));
    });
  });
  
  req.on('error', function(e) {
    console.log("ERROR:");
    console.log(e);
  });
  
  req.write(JSON.stringify(data));
  req.end();
};