import { Singularity } from './core/singularity';
import express from 'express' ; 
const app = express()

const HOST = 'localhost' ; 
const PORT = 1234 ; 

const sing = new Singularity(HOST , PORT) ; 

let i = 0 ; 

async function createQ(qname:string){
  
  try{
    const x= await sing.connect2(sing.createQueue2 ,qname) ;
    console.log(x , "queue created")

  }
  catch(err){
    console.log('err') ;
  }
}

createQ("somethingX") ; 

app.get("/", async (req ,res)=>{
  //await sing.publish("something","message",{
  //  queue: "something"
  //})
  //await sing.connect2(sing.publish , "somethingX" ,"message",{
  //  queue:"somethingX"
  //});
  await sing.publish("anything"+(++i), {
    queue: "somethingX"
  })
  res.send("hello")
})

app.listen(8181 , ()=>{
  console.log("server on 8181")
})
