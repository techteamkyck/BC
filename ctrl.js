'use strict';

const Thing = require('./thing.model');
const BlockchainService = require('../../../blockchainServices/blockchainSrvc.js');
const enrollID = require('../../../utils/enrollID')

/*
    Retrieve list of all things

    METHOD: GET
    URL : /api/v1/thing
    Response:
        [{'thing'}, {'thing'}]
*/
exports.list = function(req, res) {
    console.log("-- Query all things --")
    
    var userID = enrollID.getID(req);
    const functionName = "get_all_things"
    const args = [userID];
    const enrollmentId = userID;
    
    BlockchainService.query(functionName,args,enrollmentId).then(function(things){
        if (!things) {
            res.json([]);
        } else {
            console.log("Retrieved things from the blockchain: # " + things.length);
            res.json(things)
        }
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Retrieve thing object

    METHOD: GET
    URL: /api/v1/thing/:thingId
    Response:
        { thing }
*/
// exports.detail = function(req, res) {
//     console.log("-- Query thing --")
    
//     const functionName = "get_thing"
//     const args = [req.params.thingId];
//     const enrollmentId = enrollID.getID(req);
    
//     BlockchainService.query(functionName,args,enrollmentId).then(function(thing){
//         if (!thing) {
//             res.json([]);
//         } else {
//             console.log("Retrieved thing from the blockchain");
//             res.json(thing)
//         }
//     }).catch(function(err){
//         console.log("Error", err);
//         res.sendStatus(500);   
//     }); 
// }

/*
    Add thing object

    METHOD: POST
    URL: /api/v1/thing/
    Response:
        {  }
*/
exports.add = function(req, res) {
    console.log("-- Adding thing --")
      
    const functionName = "add_thing"
    const args = [req.body.thingId, JSON.stringify(req.body.thing)];
    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.sendStatus(200);
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}
exports.addresource = function(req, res) {
    console.log("-- Nodejs Adding resource --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "add_resource"
    const args = [req.body.owner, req.body.hash,req.body.path];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.sendStatus(200);
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

exports.getresource = function(req, res) {
    console.log("-- Nodejs Getting resource --")

    const functionName = "get_resource"
    var owner = req.param('owner')
    console.log("OWNER :: " + owner);
    
    var hash = req.param('hash')
    console.log("hash :: " + hash);
    
    const args = [owner, hash]

    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.query(functionName,args,enrollmentId).then(function(things){
        if (!things) {
            res.json([]);
        } else {
            console.log("Retrieved things from the blockchain: # " + things.length);
            res.json(things)
        }
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

exports.createBrokerageRequest = function(req, res) {
    console.log("-- Nodejs Adding resource --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "create_brokerageRequest"
    const args = [JSON.stringify(req.body)];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        var otherArray = ["item1", "item2"];
        var otherObject = { item1: "item1val", item2: "item2val" };
        var json = JSON.stringify({ 
            anObject: otherObject, 
            anArray: otherArray, 
            another: "item"
        });
        res.end(json);
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Function to Save meeting data of the application.
*/
exports.updateMeeting = function(req, res) {
    console.log("-- Nodejs Adding resource --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "update_brokerage_application"
    const args = ["MEETING", req.body.meeting, req.body.requestId];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        
        res.end(json);
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Function to Save video data of the application.
*/
exports.updateVideo = function(req, res) {
    console.log("-- Nodejs Adding resource --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "update_brokerage_application"
    const args = ["VIDEO", JSON.stringify(req.body.video), req.body.requestId];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        
        res.end(json);
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Function to Save User of the application.
*/
exports.createUser = function(req, res) {
    console.log("-- Nodejs Adding User --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "create_user"
    const args = [JSON.stringify(req.body.requestId)];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        res.end(JSON.stringify(thing));
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Function to Save User of the application.
*/
exports.updateUser = function(req, res) {
    console.log("-- Nodejs Updating User --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "update_user"
    const args = [JSON.stringify(req.body)];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        res.end(JSON.stringify(thing));
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}

/*
    Function to Get User of the application.
*/
exports.validateUser = function(req, res) {
    console.log("-- Nodejs Updating User --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "update_user"
    const args = [JSON.stringify(req.body)];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        res.end(JSON.stringify(thing));
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}


/*
    Function to Get User of the application.
*/
exports.getUser = function(req, res) {
    console.log("-- Nodejs Getting User --")
    console.log("POST BODY >>>>" + req.body);
    const functionName = "get_user"
    const args = ["STATUS", req.body.status, req.body.requestId];
    console.log("This is argument ******* " + JSON.stringify(req.body));
    const enrollmentId = enrollID.getID(req);
    BlockchainService.invoke(functionName,args,enrollmentId).then(function(thing){
        res.writeHead(200, {"Content-Type": "application/json"});
        res.end(JSON.stringify(thing));
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
}
