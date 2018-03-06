var app = new Vue({
    el: '#app',
    data: {
    //   message: 'Hello',
      instances: [{ID: '123', Status: 'runnning'}]
    },
    methods: {
      getData() {
        var route = '/check/master/cerdanyola.k8s.sandbox.nutmeg.co.uk';
        this.$http.get(route).then(response => {

            // get body data
            // self.instances = response.body.Instances;
            console.log(response.body.Instances);
            
            app.instances = response.body.Instances;
        
          }, response => {
              
          });
      }
    },
    created: function(){
        this.getData();
    },
});


var example1 = new Vue({
    el: '#example-1',
    data: {
      items: [
        { message: 'Foo' },
        { message: 'Bar' }
      ]
    }
  })