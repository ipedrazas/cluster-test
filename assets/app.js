var app = new Vue({
    el: '#app',
    data: {
    //   message: 'Hello',
      instances: [{ID: '123', Status: 'runnning'}],
      clusterName: "",
      cluster: "",
    },
    methods: {
      getData() {

        var route = '/check/master/' + this.clusterName + ".k8s.sandbox.nutmeg.co.uk";
        this.$http.get(route).then(response => {

            // get body data
            // self.instances = response.body.Instances;
            console.log(response.body.Instances);
            
            app.instances = response.body.Instances;
            app.cluster = response.body.Cluster;
        
          }, response => {
              
          });
      },
      deleteMaster(){
          // r.HandleFunc("/master/{cluster}/{all}", MastersHandler)
          var route = '/master/' + this.clusterName + '/';
        this.$http.get(route).then(response => {

            // get body data
            // self.instances = response.body.Instances;
            console.log(response.body.Instances);
            
            app.instances = response.body.Instances;
            app.cluster = response.body.Cluster;
        
          }, response => {
              
          });
      }
    },
    created: function(){
        // this.getData();
    },
});

