'use server'

import https from 'https'
import fs from 'fs'
// import fetch from 'node-fetch'

const ca = fs.readFileSync('/app/certs/ca.crt')
const cert = fs.readFileSync('/app/certs/client.crt')
const key = fs.readFileSync('/app/certs/client.key')

export const secureAgent = new https.Agent({
  ca,
  cert,
  key,
  rejectUnauthorized: true,
})
