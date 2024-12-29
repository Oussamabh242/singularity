import * as net from 'net'  ;
import { Packet, PacketType } from './packet';

export class Singularity {
  PORT   : number ;
  HOST   : string ; 
  client ?: net.Socket ; 
  constructor(host : string , port : number){
    this.PORT = port ; 
    this.HOST = host ; 
    //this.client = this.connect() ; 
  }
  private connect():net.Socket {
    const client = new net.Socket() ; 
    client.connect(this.PORT , this.HOST, ()=>{
      console.log("Socket Connected")
    }); 
    return client ; 
  }

  public connect2(cb : any , ...args : any){
    const client = new net.Socket() ; 

    client.connect(this.PORT , this.HOST,async  ()=>{
      console.log("Socket Connected")
      await cb(client , ...args) 
    }); 
    this.client = client; 
  }
  public  createQueue2(conn : net.Socket , queueName: string):Promise<any>{
    return new Promise((resolve, reject)=>{
      if (!conn.writable) {
         reject(new Error("Socket is disconnected"));
      }

      const packet = new Packet(5, "", {
        "queue": queueName,
        "content-type": "json"
      });


      conn.write(packet.encode(), (err:any) => {
        if (err) {
          reject(new Error("Failed to send packet:" +err.message));
        }
      });

      const handleData =  (ack:Buffer) => {
        const uint8array = new Uint8Array(ack);
        console.log(uint8array , ack) ;
        const decoded = Packet.decode(uint8array);
        console.log(PacketType[decoded.pType] )
        if ('ACKCREATE' != PacketType[decoded.pType]) {
          return reject(new Error("Something went wrong PACKET didnot get acknowledged"))
        }
         
        resolve(0);
      };
      conn.on("data" , handleData) ; 

    })
  }

  public  createQueue(queueName: string):Promise<any>{
    return new Promise((resolve, reject)=>{
       
      if (!this.client) {
        reject(new Error("Socket is disconnected"));
        return  ;
      }

      const packet = new Packet(5, "", {
        "queue": queueName,
        "content-type": "json"
      });


      this.client.write(packet.encode(), (err) => {
        if (err) {
          reject(new Error("Failed to send packet:" +err.message));
        }
      });

      const handleData =  (ack:Buffer) => {
        const uint8array = new Uint8Array(ack);
        console.log(uint8array , ack) ;
        const decoded = Packet.decode(uint8array);
        console.log(PacketType[decoded.pType] )
        if ('ACKCREATE' != PacketType[decoded.pType]) {
          if(!this.client){return reject("no client")}
          this.client.removeListener('data', handleData)
          return reject(new Error("Something went wrong PACKET didnot get acknowledged"))

        }
         
        resolve(0);
      };
      this.client.on("data" , handleData) ; 

    })
  }



  private sendPacket(conn : net.Socket   , payload:any , metadata:object):Promise<void>{
    return new Promise((resolve, reject)=>{
      if(!conn){
        reject("no client") ;
        return 
      }
      if(!conn.writable){
        reject(new Error("Socket is closed")); 
      }
      const packet = new Packet(1 ,payload ,metadata); 
      conn.write(packet.encode(), (err) => {
        if (err) {
          return reject(new Error("Failed to send publish packet: " + err.message));
        }
        console.log("Message published");
        resolve(); // Successfully published the message
      });
    });
      
  }

  private listenToQueue(conn : net.Socket, queueName:any , fn:(arg:any)=>void | Promise<void> ){
    if(!conn){
      throw new Error("Socket disconnected") ; 
      return
    }
    if(!conn.readable){
      new Error("Socket is closed"); 
    }
    const packet = new Packet(PacketType.SUBSCRIBE ,"" ,{
      queue : queueName
    }); 
    conn.write(packet.encode(),(err)=>{
      console.error("something went wrong when subscribing", err);
    });
    conn.on("data" ,(data:Buffer)=>{
      const arr = new Uint8Array(data); 
      //console.log(arr)

      if(arr[0]!= PacketType.JOB){
        console.log("not JOB")
      }
      else{
        const encoded =  Packet.decode(arr) ;

        const res = fn(encoded.payload)

        if(res && typeof res.then == 'function'){
          res.then(()=>{
            conn.write("ok", (err)=>{console.error(err)}) ; 
          }); 
        }
        else{
          conn.write("ok" , (err)=>{console.log(err)})
        }

        
        
      }

    })
      
    
  }

  public publish(payload:any , metadata:any){
    this.connect2(this.sendPacket ,payload , metadata)
  }
  public subscribe(queueName:string , job :(arg:any)=>void){
    this.connect2(this.listenToQueue,queueName ,job);
  }



  public closeConnection(): void {
    if (this.client) {
      this.client.end(() => {
        console.log("Connection closed");
      });
    }
  }
}


function waitForRandomTime(arg : any): Promise<void> {
  return new Promise(resolve => {
    const randomTime = Math.floor(Math.random() * 4000) + 1000; 

    console.log(`working on ${arg} takes ${randomTime} milliseconds...`);

    setTimeout(() => {
      console.log(`Job Done \n--------------------`);
      resolve();
    }, randomTime);
  });
}


async function main(){
  const sing = new Singularity("localhost" , 1234) ; 
  sing.subscribe("somethingX",waitForRandomTime)
  //try{
  //     sing.connect2(sing.createQueue2 ,"somethingelse "); 
  //}
  //catch(err){
  //  console.error( err) 
  //}
  //
  //sing.connect2(sing.publish , "something","message" ,{
  //  queue :"something" 
  //})
}


main();
