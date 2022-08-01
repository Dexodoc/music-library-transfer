"use strict";

const fs = require("fs");
const jwt = require("jsonwebtoken");

const privateKey = fs.readFileSync("AuthKey.p8").toString();
const teamId = "3S59UFWP8Y";
const keyId = "TRUC777DZL";

const jwtToken = jwt.sign({}, privateKey, {
  algorithm: "ES256",
  expiresIn: "180d",
  issuer: teamId,
  header: {
    alg: "ES256",
    kid: keyId
  }
});

console.log(jwtToken);