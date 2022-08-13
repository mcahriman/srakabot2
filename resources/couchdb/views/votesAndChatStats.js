function (keys, values, rereduce) {
    let sumKarma = function(a,b) {
      return a+b.karma;
    };
    let addMessageCount = function(a,b) {
      return a+b.messageCount;
    };
    let addVote = function(a,b) {
      return a+b.vote;
    };
    
    if(rereduce) {
      return {
        karma: values.reduce(sumKarma, 0),
        messageCount: values.reduce(addMessageCount, 0),
        user: values[0].user
      };
    }
    return {
      karma: values.reduce(addVote ,0),
      messageCount: values.reduce(addMessageCount,0),
      user: values[0].user
    };
  }