import slo from 'k6/x/slo-generate'
import http from 'k6/http'
import { group } from 'k6'


const config = {
  percentage: "99.9",
}

const responseData = []

export const options = {
  iterations: 50
}

// change this to be more specific 
// change this to allow you to set p99 or p95
function handleResponses(responseTime) { 
  
  responseData.push(responseTime)

  // & responseData.length >= 50
  if (responseData.length == options.iterations -1 ) { 

    
    // actually generate the threashold 
    const slo_data = slo.generate(responseData, config.percentage)

    console.log("THREASHOLD")
    console.log(slo_data[0])

    console.log("SLO")
    console.log(slo_data[1])
    
  } 

}
 
export default function () {
    // get data from the actual run of a service 
    let response
    let name

    group('First API call', function () {
      name = "GET /"
      response = http.get('http://test.k6.io')
      handleResponses(response.timings.duration)

    })


  }
  