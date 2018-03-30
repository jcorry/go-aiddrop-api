const siege = require('siege')

siege()
  .on(9090)
  .concurrent(25)
  .get('/cases?maxLongitude=-84.285&minLongitude=-84.290&minLatitude=33.8628&maxLatitude=33.8655')
  .attack()