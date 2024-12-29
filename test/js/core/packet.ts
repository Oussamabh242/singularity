
export enum PacketType {
  PUBLISH = 1,         // Publish
  ACKRECIVE = 2,       // Acknowledge Publishing
  SUBSCRIBE = 3,       // Subscribe to Queue
  ACKSUBSCRIBE = 4,    // Acknowledge Subscribe
  CREATEQUEUE = 5,     // Create Queue
  ACKCREATE = 6,       // Acknowledge Create
  JOB = 7,             // Job
  ACKJOB = 8,          // Acknowledge Job
  PING = 100           // Ping
}

export class Packet {
  pType    : number ; 
  payload  : string ; 
  metadata : string ; 
  rLength  : number ;
  mLength  : number ; 
  pLength  : number ;

  constructor (pType : number , payload : any , metadata : object){
    this.metadata = this.stringify(metadata) ;
    this.payload =  this.stringify(payload)   ;
    this.pType   =  pType;
    this.pLength = this.payload.length ; 
    this.mLength = this.metadata.length ;
    this.rLength = this.mLength+this.pLength+4; 
  }

  public encode():Uint8Array{
    const eMeta = this.strToBytes(this.metadata) ;
    const ePayload = this.strToBytes(this.payload) ; 
    const epType  = this.intToBytes(this.pType,1)  ;
    const erLength = this.intToBytes(this.rLength,4) ; 
    const emLength = this.intToBytes(this.mLength,2) ; 
    const epLength = this.intToBytes(this.pLength,2) ;
    let packet = new Uint8Array(this.rLength +5) ; 
    let all = [epType , erLength ,emLength , eMeta , epLength, ePayload]; 
    let idx = 0 ; 
    for(let arr of all){
      packet.set(arr ,idx) ;
      idx+=arr.length ; 
    }
    
    return packet

  }

  static decode(packet : Uint8Array):Packet{
    const ptype = packet[0];
    const rem_len = new DataView(packet.buffer).getUint32(1, false); 
    if (rem_len == 4) {
      return new Packet(ptype ,"" , {}) ;
    }
    const meta_len = new DataView(packet.buffer).getUint16(5, false); 
        
    let metadata:any = {};
    let nextByte = 7;
        
    if (meta_len > 0) {
      metadata = String.fromCharCode.apply(null, [...packet.slice(7, 7 + meta_len)]);
      nextByte = 7 + meta_len;
      metadata = JSON.parse(metadata)
    }

        
        // Read payload length (2 bytes, unsigned short)
    const payload_len = new DataView(packet.buffer).getUint16(nextByte, false); // false = big-endian
    let payload = "";
        
    if (payload_len > 0) {
      payload = String.fromCharCode.apply(null, [...packet.slice(nextByte + 2, nextByte + 2 + payload_len)]);
    }
        
    //console.log(ptype, metadata, payload);
        
    return new Packet(ptype, payload, metadata);
        
  }
  
  private stringify(thing : any){
    if(typeof thing === 'object' && thing.length === undefined){
      return JSON.stringify(thing) ;
    }
    return thing+"" ; 
  }
  
  private intToBytes(int:number , length:number) : Uint8Array {
    const buffer = new ArrayBuffer(length);  
    const view = new DataView(buffer);

    if(length == 1){
      view.setUint8(0,int) ; 
    }
    else if(length == 2) {
        view.setInt16(0, int, false);  
    }else {
        view.setInt32(0, int, false);
    }

    return new Uint8Array(buffer); 
  }
  private strToBytes(str : string):Uint8Array{
    const encoder = new TextEncoder() ; 
    return encoder.encode(str) ; 
  }
  
}

//const array = new Uint8Array([7, 0, 0, 0, 12, 0, 0, 0, 8, 97, 110, 121, 116, 104, 105, 110, 103]);
//const x = Packet.decode(array)
//console.log(x.pType , x.payload , x.metadata)

//const y = new Uint8Array([1,0,0, 18, 52, 0, 17, 123, 34, 113, 117, 101, 117, 101, 34, 58, 32, 34, 116, 104, 97, 116, 34, 125, 0, 12, 104, 101, 108, 108, 111, 32, 116, 104, 101, 114, 101, 48])
//let x = Packet.decode(y)
//console.log(x.mLength)
//console.log(Packet.decode())
//function stringify(thing : any):string{
//    if(typeof thing === 'object' && thing.length === undefined){
//      return JSON.stringify(this.payload) ;
//    }
//    return thing+"" ; 
//  }

//console.log(Packet.stringify())

//const packet = new Packet(1, '', {
//			Queue:       "thing"
//		});
//let x = packet.encode()
//
//console.log(x) ; 

//console.log("pType:", packet.pType);  // Expected output: 1
//console.log("payload:", packet.payload);  // Expected output: 'Hello, World!'
//console.log("metadata:", packet.metadata);  // Expected output: '{"source":"client","destination":"server"}'
//console.log("pLength:", packet.pLength);  // Expected output: length of 'Hello, World!' = 13
//console.log("mLength:", packet.mLength);  // Expected output: length of '{"source":"client","destination":"server"}' = 43
//console.log("rLength:", packet.rLength);  // Expected output: mLength + pLength + 4 = 43 + 13 + 4 = 60

