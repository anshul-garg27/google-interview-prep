import * as jspb from 'google-protobuf'



export class Events extends jspb.Message {
  getHeader(): Header | undefined;
  setHeader(value?: Header): Events;
  hasHeader(): boolean;
  clearHeader(): Events;

  getEventsList(): Array<Event>;
  setEventsList(value: Array<Event>): Events;
  clearEventsList(): Events;
  addEvents(value?: Event, index?: number): Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Events.AsObject;
  static toObject(includeInstance: boolean, msg: Events): Events.AsObject;
  static serializeBinaryToWriter(message: Events, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Events;
  static deserializeBinaryFromReader(message: Events, reader: jspb.BinaryReader): Events;
}

export namespace Events {
  export type AsObject = {
    header?: Header.AsObject,
    eventsList: Array<Event.AsObject>,
  }
}

export class Header extends jspb.Message {
  getSessionid(): string;
  setSessionid(value: string): Header;

  getBbdeviceid(): number;
  setBbdeviceid(value: number): Header;

  getUserid(): number;
  setUserid(value: number): Header;

  getDeviceid(): string;
  setDeviceid(value: string): Header;

  getAndroidadvertisingid(): string;
  setAndroidadvertisingid(value: string): Header;

  getClientid(): string;
  setClientid(value: string): Header;

  getChannel(): string;
  setChannel(value: string): Header;

  getOs(): string;
  setOs(value: string): Header;

  getClienttype(): string;
  setClienttype(value: string): Header;

  getApplanguage(): string;
  setApplanguage(value: string): Header;

  getAppversion(): string;
  setAppversion(value: string): Header;

  getAppversioncode(): string;
  setAppversioncode(value: string): Header;

  getCurrenturl(): string;
  setCurrenturl(value: string): Header;

  getMerchantid(): string;
  setMerchantid(value: string): Header;

  getUtmreferrer(): string;
  setUtmreferrer(value: string): Header;

  getUtmplatform(): string;
  setUtmplatform(value: string): Header;

  getUtmsource(): string;
  setUtmsource(value: string): Header;

  getUtmmedium(): string;
  setUtmmedium(value: string): Header;

  getUtmcampaign(): string;
  setUtmcampaign(value: string): Header;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Header.AsObject;
  static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
  static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Header;
  static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
}

export namespace Header {
  export type AsObject = {
    sessionid: string,
    bbdeviceid: number,
    userid: number,
    deviceid: string,
    androidadvertisingid: string,
    clientid: string,
    channel: string,
    os: string,
    clienttype: string,
    applanguage: string,
    appversion: string,
    appversioncode: string,
    currenturl: string,
    merchantid: string,
    utmreferrer: string,
    utmplatform: string,
    utmsource: string,
    utmmedium: string,
    utmcampaign: string,
  }
}

export class Event extends jspb.Message {
  getId(): string;
  setId(value: string): Event;

  getTimestamp(): string;
  setTimestamp(value: string): Event;

  getLaunchevent(): LaunchEvent | undefined;
  setLaunchevent(value?: LaunchEvent): Event;
  hasLaunchevent(): boolean;
  clearLaunchevent(): Event;

  getPageopenedevent(): PageOpenedEvent | undefined;
  setPageopenedevent(value?: PageOpenedEvent): Event;
  hasPageopenedevent(): boolean;
  clearPageopenedevent(): Event;

  getPageloadedevent(): PageLoadedEvent | undefined;
  setPageloadedevent(value?: PageLoadedEvent): Event;
  hasPageloadedevent(): boolean;
  clearPageloadedevent(): Event;

  getPageloadfailedevent(): PageLoadFailedEvent | undefined;
  setPageloadfailedevent(value?: PageLoadFailedEvent): Event;
  hasPageloadfailedevent(): boolean;
  clearPageloadfailedevent(): Event;

  getWidgetviewevent(): WidgetViewEvent | undefined;
  setWidgetviewevent(value?: WidgetViewEvent): Event;
  hasWidgetviewevent(): boolean;
  clearWidgetviewevent(): Event;

  getWidgetctaclickedevent(): WidgetCtaClickedEvent | undefined;
  setWidgetctaclickedevent(value?: WidgetCtaClickedEvent): Event;
  hasWidgetctaclickedevent(): boolean;
  clearWidgetctaclickedevent(): Event;

  getWidgetelementviewevent(): WidgetElementViewEvent | undefined;
  setWidgetelementviewevent(value?: WidgetElementViewEvent): Event;
  hasWidgetelementviewevent(): boolean;
  clearWidgetelementviewevent(): Event;

  getWidgetelementclickedevent(): WidgetElementClickedEvent | undefined;
  setWidgetelementclickedevent(value?: WidgetElementClickedEvent): Event;
  hasWidgetelementclickedevent(): boolean;
  clearWidgetelementclickedevent(): Event;

  getEnterpwaproductpageevent(): EnterPWAProductPageEvent | undefined;
  setEnterpwaproductpageevent(value?: EnterPWAProductPageEvent): Event;
  hasEnterpwaproductpageevent(): boolean;
  clearEnterpwaproductpageevent(): Event;

  getStreamenterevent(): StreamEnterEvent | undefined;
  setStreamenterevent(value?: StreamEnterEvent): Event;
  hasStreamenterevent(): boolean;
  clearStreamenterevent(): Event;

  getChooseproductevent(): ChooseProductEvent | undefined;
  setChooseproductevent(value?: ChooseProductEvent): Event;
  hasChooseproductevent(): boolean;
  clearChooseproductevent(): Event;

  getPhoneverificationinitiateevent(): PhoneVerificationInitiateEvent | undefined;
  setPhoneverificationinitiateevent(value?: PhoneVerificationInitiateEvent): Event;
  hasPhoneverificationinitiateevent(): boolean;
  clearPhoneverificationinitiateevent(): Event;

  getPhonenumberentered(): PhoneNumberEntered | undefined;
  setPhonenumberentered(value?: PhoneNumberEntered): Event;
  hasPhonenumberentered(): boolean;
  clearPhonenumberentered(): Event;

  getOtpverifiedevent(): OTPVerifiedEvent | undefined;
  setOtpverifiedevent(value?: OTPVerifiedEvent): Event;
  hasOtpverifiedevent(): boolean;
  clearOtpverifiedevent(): Event;

  getAddtocartevent(): AddToCartEvent | undefined;
  setAddtocartevent(value?: AddToCartEvent): Event;
  hasAddtocartevent(): boolean;
  clearAddtocartevent(): Event;

  getGotopaymentsevent(): GoToPaymentsEvent | undefined;
  setGotopaymentsevent(value?: GoToPaymentsEvent): Event;
  hasGotopaymentsevent(): boolean;
  clearGotopaymentsevent(): Event;

  getInitiatepurchaseevent(): InitiatePurchaseEvent | undefined;
  setInitiatepurchaseevent(value?: InitiatePurchaseEvent): Event;
  hasInitiatepurchaseevent(): boolean;
  clearInitiatepurchaseevent(): Event;

  getCompletepurchaseevent(): CompletePurchaseEvent | undefined;
  setCompletepurchaseevent(value?: CompletePurchaseEvent): Event;
  hasCompletepurchaseevent(): boolean;
  clearCompletepurchaseevent(): Event;

  getEventofCase(): Event.EventofCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    id: string,
    timestamp: string,
    launchevent?: LaunchEvent.AsObject,
    pageopenedevent?: PageOpenedEvent.AsObject,
    pageloadedevent?: PageLoadedEvent.AsObject,
    pageloadfailedevent?: PageLoadFailedEvent.AsObject,
    widgetviewevent?: WidgetViewEvent.AsObject,
    widgetctaclickedevent?: WidgetCtaClickedEvent.AsObject,
    widgetelementviewevent?: WidgetElementViewEvent.AsObject,
    widgetelementclickedevent?: WidgetElementClickedEvent.AsObject,
    enterpwaproductpageevent?: EnterPWAProductPageEvent.AsObject,
    streamenterevent?: StreamEnterEvent.AsObject,
    chooseproductevent?: ChooseProductEvent.AsObject,
    phoneverificationinitiateevent?: PhoneVerificationInitiateEvent.AsObject,
    phonenumberentered?: PhoneNumberEntered.AsObject,
    otpverifiedevent?: OTPVerifiedEvent.AsObject,
    addtocartevent?: AddToCartEvent.AsObject,
    gotopaymentsevent?: GoToPaymentsEvent.AsObject,
    initiatepurchaseevent?: InitiatePurchaseEvent.AsObject,
    completepurchaseevent?: CompletePurchaseEvent.AsObject,
  }

  export enum EventofCase { 
    EVENTOF_NOT_SET = 0,
    LAUNCHEVENT = 9,
    PAGEOPENEDEVENT = 10,
    PAGELOADEDEVENT = 11,
    PAGELOADFAILEDEVENT = 12,
    WIDGETVIEWEVENT = 13,
    WIDGETCTACLICKEDEVENT = 14,
    WIDGETELEMENTVIEWEVENT = 15,
    WIDGETELEMENTCLICKEDEVENT = 16,
    ENTERPWAPRODUCTPAGEEVENT = 17,
    STREAMENTEREVENT = 18,
    CHOOSEPRODUCTEVENT = 19,
    PHONEVERIFICATIONINITIATEEVENT = 20,
    PHONENUMBERENTERED = 21,
    OTPVERIFIEDEVENT = 22,
    ADDTOCARTEVENT = 23,
    GOTOPAYMENTSEVENT = 24,
    INITIATEPURCHASEEVENT = 25,
    COMPLETEPURCHASEEVENT = 26,
  }
}

export class RecoInfo extends jspb.Message {
  getModel(): string;
  setModel(value: string): RecoInfo;

  getPersonatype(): string;
  setPersonatype(value: string): RecoInfo;

  getPersonavalue(): string;
  setPersonavalue(value: string): RecoInfo;

  getBucketlabel(): string;
  setBucketlabel(value: string): RecoInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RecoInfo.AsObject;
  static toObject(includeInstance: boolean, msg: RecoInfo): RecoInfo.AsObject;
  static serializeBinaryToWriter(message: RecoInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RecoInfo;
  static deserializeBinaryFromReader(message: RecoInfo, reader: jspb.BinaryReader): RecoInfo;
}

export namespace RecoInfo {
  export type AsObject = {
    model: string,
    personatype: string,
    personavalue: string,
    bucketlabel: string,
  }
}

export class LaunchEvent extends jspb.Message {
  getDynamiclink(): string;
  setDynamiclink(value: string): LaunchEvent;

  getApplaunchcount(): number;
  setApplaunchcount(value: number): LaunchEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LaunchEvent.AsObject;
  static toObject(includeInstance: boolean, msg: LaunchEvent): LaunchEvent.AsObject;
  static serializeBinaryToWriter(message: LaunchEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LaunchEvent;
  static deserializeBinaryFromReader(message: LaunchEvent, reader: jspb.BinaryReader): LaunchEvent;
}

export namespace LaunchEvent {
  export type AsObject = {
    dynamiclink: string,
    applaunchcount: number,
  }
}

export class PageOpenedEvent extends jspb.Message {
  getPagequeryid(): string;
  setPagequeryid(value: string): PageOpenedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PageOpenedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PageOpenedEvent): PageOpenedEvent.AsObject;
  static serializeBinaryToWriter(message: PageOpenedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PageOpenedEvent;
  static deserializeBinaryFromReader(message: PageOpenedEvent, reader: jspb.BinaryReader): PageOpenedEvent;
}

export namespace PageOpenedEvent {
  export type AsObject = {
    pagequeryid: string,
  }
}

export class PageLoadedEvent extends jspb.Message {
  getPagequeryid(): string;
  setPagequeryid(value: string): PageLoadedEvent;

  getBackendvariant(): string;
  setBackendvariant(value: string): PageLoadedEvent;

  getBackendsourcetype(): string;
  setBackendsourcetype(value: string): PageLoadedEvent;

  getBackendsource(): string;
  setBackendsource(value: string): PageLoadedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PageLoadedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PageLoadedEvent): PageLoadedEvent.AsObject;
  static serializeBinaryToWriter(message: PageLoadedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PageLoadedEvent;
  static deserializeBinaryFromReader(message: PageLoadedEvent, reader: jspb.BinaryReader): PageLoadedEvent;
}

export namespace PageLoadedEvent {
  export type AsObject = {
    pagequeryid: string,
    backendvariant: string,
    backendsourcetype: string,
    backendsource: string,
  }
}

export class PageLoadFailedEvent extends jspb.Message {
  getPagequeryid(): string;
  setPagequeryid(value: string): PageLoadFailedEvent;

  getReason(): string;
  setReason(value: string): PageLoadFailedEvent;

  getErrormessageshown(): string;
  setErrormessageshown(value: string): PageLoadFailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PageLoadFailedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PageLoadFailedEvent): PageLoadFailedEvent.AsObject;
  static serializeBinaryToWriter(message: PageLoadFailedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PageLoadFailedEvent;
  static deserializeBinaryFromReader(message: PageLoadFailedEvent, reader: jspb.BinaryReader): PageLoadFailedEvent;
}

export namespace PageLoadFailedEvent {
  export type AsObject = {
    pagequeryid: string,
    reason: string,
    errormessageshown: string,
  }
}

export class WidgetViewEvent extends jspb.Message {
  getVariantname(): string;
  setVariantname(value: string): WidgetViewEvent;

  getWidgettype(): string;
  setWidgettype(value: string): WidgetViewEvent;

  getElementcount(): number;
  setElementcount(value: number): WidgetViewEvent;

  getPageslot(): number;
  setPageslot(value: number): WidgetViewEvent;

  getWidgetid(): string;
  setWidgetid(value: string): WidgetViewEvent;

  getRecoinfo(): RecoInfo | undefined;
  setRecoinfo(value?: RecoInfo): WidgetViewEvent;
  hasRecoinfo(): boolean;
  clearRecoinfo(): WidgetViewEvent;

  getTitle(): string;
  setTitle(value: string): WidgetViewEvent;

  getSubtitle(): string;
  setSubtitle(value: string): WidgetViewEvent;

  getLabel(): string;
  setLabel(value: string): WidgetViewEvent;

  getElementtype(): string;
  setElementtype(value: string): WidgetViewEvent;

  getHascta(): boolean;
  setHascta(value: boolean): WidgetViewEvent;

  getSource(): string;
  setSource(value: string): WidgetViewEvent;

  getSourceidentifier(): string;
  setSourceidentifier(value: string): WidgetViewEvent;

  getHeadervariantname(): string;
  setHeadervariantname(value: string): WidgetViewEvent;

  getFootervariantname(): string;
  setFootervariantname(value: string): WidgetViewEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WidgetViewEvent.AsObject;
  static toObject(includeInstance: boolean, msg: WidgetViewEvent): WidgetViewEvent.AsObject;
  static serializeBinaryToWriter(message: WidgetViewEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WidgetViewEvent;
  static deserializeBinaryFromReader(message: WidgetViewEvent, reader: jspb.BinaryReader): WidgetViewEvent;
}

export namespace WidgetViewEvent {
  export type AsObject = {
    variantname: string,
    widgettype: string,
    elementcount: number,
    pageslot: number,
    widgetid: string,
    recoinfo?: RecoInfo.AsObject,
    title: string,
    subtitle: string,
    label: string,
    elementtype: string,
    hascta: boolean,
    source: string,
    sourceidentifier: string,
    headervariantname: string,
    footervariantname: string,
  }
}

export class WidgetCtaClickedEvent extends jspb.Message {
  getVariantname(): string;
  setVariantname(value: string): WidgetCtaClickedEvent;

  getWidgettype(): string;
  setWidgettype(value: string): WidgetCtaClickedEvent;

  getElementcount(): number;
  setElementcount(value: number): WidgetCtaClickedEvent;

  getPageslot(): number;
  setPageslot(value: number): WidgetCtaClickedEvent;

  getWidgetid(): string;
  setWidgetid(value: string): WidgetCtaClickedEvent;

  getRecoinfo(): RecoInfo | undefined;
  setRecoinfo(value?: RecoInfo): WidgetCtaClickedEvent;
  hasRecoinfo(): boolean;
  clearRecoinfo(): WidgetCtaClickedEvent;

  getTitle(): string;
  setTitle(value: string): WidgetCtaClickedEvent;

  getSubtitle(): string;
  setSubtitle(value: string): WidgetCtaClickedEvent;

  getLabel(): string;
  setLabel(value: string): WidgetCtaClickedEvent;

  getActionlink(): string;
  setActionlink(value: string): WidgetCtaClickedEvent;

  getElementtype(): string;
  setElementtype(value: string): WidgetCtaClickedEvent;

  getCtatext(): string;
  setCtatext(value: string): WidgetCtaClickedEvent;

  getSource(): string;
  setSource(value: string): WidgetCtaClickedEvent;

  getSourceidentifier(): string;
  setSourceidentifier(value: string): WidgetCtaClickedEvent;

  getHeadervariantname(): string;
  setHeadervariantname(value: string): WidgetCtaClickedEvent;

  getFootervariantname(): string;
  setFootervariantname(value: string): WidgetCtaClickedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WidgetCtaClickedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: WidgetCtaClickedEvent): WidgetCtaClickedEvent.AsObject;
  static serializeBinaryToWriter(message: WidgetCtaClickedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WidgetCtaClickedEvent;
  static deserializeBinaryFromReader(message: WidgetCtaClickedEvent, reader: jspb.BinaryReader): WidgetCtaClickedEvent;
}

export namespace WidgetCtaClickedEvent {
  export type AsObject = {
    variantname: string,
    widgettype: string,
    elementcount: number,
    pageslot: number,
    widgetid: string,
    recoinfo?: RecoInfo.AsObject,
    title: string,
    subtitle: string,
    label: string,
    actionlink: string,
    elementtype: string,
    ctatext: string,
    source: string,
    sourceidentifier: string,
    headervariantname: string,
    footervariantname: string,
  }
}

export class WidgetElementViewEvent extends jspb.Message {
  getVariantname(): string;
  setVariantname(value: string): WidgetElementViewEvent;

  getWidgettype(): string;
  setWidgettype(value: string): WidgetElementViewEvent;

  getElementcount(): number;
  setElementcount(value: number): WidgetElementViewEvent;

  getPageslot(): number;
  setPageslot(value: number): WidgetElementViewEvent;

  getWidgetid(): string;
  setWidgetid(value: string): WidgetElementViewEvent;

  getPosition(): number;
  setPosition(value: number): WidgetElementViewEvent;

  getRecoinfo(): RecoInfo | undefined;
  setRecoinfo(value?: RecoInfo): WidgetElementViewEvent;
  hasRecoinfo(): boolean;
  clearRecoinfo(): WidgetElementViewEvent;

  getTitle(): string;
  setTitle(value: string): WidgetElementViewEvent;

  getSubtitle(): string;
  setSubtitle(value: string): WidgetElementViewEvent;

  getLabel(): string;
  setLabel(value: string): WidgetElementViewEvent;

  getElementid(): string;
  setElementid(value: string): WidgetElementViewEvent;

  getElementtype(): string;
  setElementtype(value: string): WidgetElementViewEvent;

  getFomotext(): string;
  setFomotext(value: string): WidgetElementViewEvent;

  getPrice(): number;
  setPrice(value: number): WidgetElementViewEvent;

  getSource(): string;
  setSource(value: string): WidgetElementViewEvent;

  getSourceidentifier(): string;
  setSourceidentifier(value: string): WidgetElementViewEvent;

  getHeadervariantname(): string;
  setHeadervariantname(value: string): WidgetElementViewEvent;

  getFootervariantname(): string;
  setFootervariantname(value: string): WidgetElementViewEvent;

  getProducttype(): string;
  setProducttype(value: string): WidgetElementViewEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WidgetElementViewEvent.AsObject;
  static toObject(includeInstance: boolean, msg: WidgetElementViewEvent): WidgetElementViewEvent.AsObject;
  static serializeBinaryToWriter(message: WidgetElementViewEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WidgetElementViewEvent;
  static deserializeBinaryFromReader(message: WidgetElementViewEvent, reader: jspb.BinaryReader): WidgetElementViewEvent;
}

export namespace WidgetElementViewEvent {
  export type AsObject = {
    variantname: string,
    widgettype: string,
    elementcount: number,
    pageslot: number,
    widgetid: string,
    position: number,
    recoinfo?: RecoInfo.AsObject,
    title: string,
    subtitle: string,
    label: string,
    elementid: string,
    elementtype: string,
    fomotext: string,
    price: number,
    source: string,
    sourceidentifier: string,
    headervariantname: string,
    footervariantname: string,
    producttype: string,
  }
}

export class WidgetElementClickedEvent extends jspb.Message {
  getVariantname(): string;
  setVariantname(value: string): WidgetElementClickedEvent;

  getWidgettype(): string;
  setWidgettype(value: string): WidgetElementClickedEvent;

  getElementcount(): number;
  setElementcount(value: number): WidgetElementClickedEvent;

  getPageslot(): number;
  setPageslot(value: number): WidgetElementClickedEvent;

  getWidgetid(): string;
  setWidgetid(value: string): WidgetElementClickedEvent;

  getPosition(): number;
  setPosition(value: number): WidgetElementClickedEvent;

  getRecoinfo(): RecoInfo | undefined;
  setRecoinfo(value?: RecoInfo): WidgetElementClickedEvent;
  hasRecoinfo(): boolean;
  clearRecoinfo(): WidgetElementClickedEvent;

  getTitle(): string;
  setTitle(value: string): WidgetElementClickedEvent;

  getSubtitle(): string;
  setSubtitle(value: string): WidgetElementClickedEvent;

  getLabel(): string;
  setLabel(value: string): WidgetElementClickedEvent;

  getElementid(): string;
  setElementid(value: string): WidgetElementClickedEvent;

  getElementtype(): string;
  setElementtype(value: string): WidgetElementClickedEvent;

  getFomotext(): string;
  setFomotext(value: string): WidgetElementClickedEvent;

  getPrice(): number;
  setPrice(value: number): WidgetElementClickedEvent;

  getSource(): string;
  setSource(value: string): WidgetElementClickedEvent;

  getSourceidentifier(): string;
  setSourceidentifier(value: string): WidgetElementClickedEvent;

  getHeadervariantname(): string;
  setHeadervariantname(value: string): WidgetElementClickedEvent;

  getFootervariantname(): string;
  setFootervariantname(value: string): WidgetElementClickedEvent;

  getProducttype(): string;
  setProducttype(value: string): WidgetElementClickedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WidgetElementClickedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: WidgetElementClickedEvent): WidgetElementClickedEvent.AsObject;
  static serializeBinaryToWriter(message: WidgetElementClickedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WidgetElementClickedEvent;
  static deserializeBinaryFromReader(message: WidgetElementClickedEvent, reader: jspb.BinaryReader): WidgetElementClickedEvent;
}

export namespace WidgetElementClickedEvent {
  export type AsObject = {
    variantname: string,
    widgettype: string,
    elementcount: number,
    pageslot: number,
    widgetid: string,
    position: number,
    recoinfo?: RecoInfo.AsObject,
    title: string,
    subtitle: string,
    label: string,
    elementid: string,
    elementtype: string,
    fomotext: string,
    price: number,
    source: string,
    sourceidentifier: string,
    headervariantname: string,
    footervariantname: string,
    producttype: string,
  }
}

export class EnterPWAProductPageEvent extends jspb.Message {
  getProductid(): number;
  setProductid(value: number): EnterPWAProductPageEvent;

  getProducttype(): string;
  setProducttype(value: string): EnterPWAProductPageEvent;

  getHcpid(): string;
  setHcpid(value: string): EnterPWAProductPageEvent;

  getAffiliatechannel(): string;
  setAffiliatechannel(value: string): EnterPWAProductPageEvent;

  getProductname(): string;
  setProductname(value: string): EnterPWAProductPageEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnterPWAProductPageEvent.AsObject;
  static toObject(includeInstance: boolean, msg: EnterPWAProductPageEvent): EnterPWAProductPageEvent.AsObject;
  static serializeBinaryToWriter(message: EnterPWAProductPageEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnterPWAProductPageEvent;
  static deserializeBinaryFromReader(message: EnterPWAProductPageEvent, reader: jspb.BinaryReader): EnterPWAProductPageEvent;
}

export namespace EnterPWAProductPageEvent {
  export type AsObject = {
    productid: number,
    producttype: string,
    hcpid: string,
    affiliatechannel: string,
    productname: string,
  }
}

export class StreamEnterEvent extends jspb.Message {
  getStreamid(): number;
  setStreamid(value: number): StreamEnterEvent;

  getStreamtitle(): string;
  setStreamtitle(value: string): StreamEnterEvent;

  getStreamdescription(): string;
  setStreamdescription(value: string): StreamEnterEvent;

  getIslive(): boolean;
  setIslive(value: boolean): StreamEnterEvent;

  getProductidsList(): Array<number>;
  setProductidsList(value: Array<number>): StreamEnterEvent;
  clearProductidsList(): StreamEnterEvent;
  addProductids(value: number, index?: number): StreamEnterEvent;

  getSourceurl(): string;
  setSourceurl(value: string): StreamEnterEvent;

  getSourcecomponent(): string;
  setSourcecomponent(value: string): StreamEnterEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamEnterEvent.AsObject;
  static toObject(includeInstance: boolean, msg: StreamEnterEvent): StreamEnterEvent.AsObject;
  static serializeBinaryToWriter(message: StreamEnterEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamEnterEvent;
  static deserializeBinaryFromReader(message: StreamEnterEvent, reader: jspb.BinaryReader): StreamEnterEvent;
}

export namespace StreamEnterEvent {
  export type AsObject = {
    streamid: number,
    streamtitle: string,
    streamdescription: string,
    islive: boolean,
    productidsList: Array<number>,
    sourceurl: string,
    sourcecomponent: string,
  }
}

export class ChooseProductEvent extends jspb.Message {
  getStreamid(): number;
  setStreamid(value: number): ChooseProductEvent;

  getStreamtitle(): string;
  setStreamtitle(value: string): ChooseProductEvent;

  getStreamdescription(): string;
  setStreamdescription(value: string): ChooseProductEvent;

  getIslive(): boolean;
  setIslive(value: boolean): ChooseProductEvent;

  getProductidsList(): Array<number>;
  setProductidsList(value: Array<number>): ChooseProductEvent;
  clearProductidsList(): ChooseProductEvent;
  addProductids(value: number, index?: number): ChooseProductEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseProductEvent.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseProductEvent): ChooseProductEvent.AsObject;
  static serializeBinaryToWriter(message: ChooseProductEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseProductEvent;
  static deserializeBinaryFromReader(message: ChooseProductEvent, reader: jspb.BinaryReader): ChooseProductEvent;
}

export namespace ChooseProductEvent {
  export type AsObject = {
    streamid: number,
    streamtitle: string,
    streamdescription: string,
    islive: boolean,
    productidsList: Array<number>,
  }
}

export class PhoneVerificationInitiateEvent extends jspb.Message {
  getSourceaction(): string;
  setSourceaction(value: string): PhoneVerificationInitiateEvent;

  getLoginflow(): string;
  setLoginflow(value: string): PhoneVerificationInitiateEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PhoneVerificationInitiateEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PhoneVerificationInitiateEvent): PhoneVerificationInitiateEvent.AsObject;
  static serializeBinaryToWriter(message: PhoneVerificationInitiateEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PhoneVerificationInitiateEvent;
  static deserializeBinaryFromReader(message: PhoneVerificationInitiateEvent, reader: jspb.BinaryReader): PhoneVerificationInitiateEvent;
}

export namespace PhoneVerificationInitiateEvent {
  export type AsObject = {
    sourceaction: string,
    loginflow: string,
  }
}

export class PhoneNumberEntered extends jspb.Message {
  getIsnewuser(): boolean;
  setIsnewuser(value: boolean): PhoneNumberEntered;

  getPhonenumber(): number;
  setPhonenumber(value: number): PhoneNumberEntered;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PhoneNumberEntered.AsObject;
  static toObject(includeInstance: boolean, msg: PhoneNumberEntered): PhoneNumberEntered.AsObject;
  static serializeBinaryToWriter(message: PhoneNumberEntered, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PhoneNumberEntered;
  static deserializeBinaryFromReader(message: PhoneNumberEntered, reader: jspb.BinaryReader): PhoneNumberEntered;
}

export namespace PhoneNumberEntered {
  export type AsObject = {
    isnewuser: boolean,
    phonenumber: number,
  }
}

export class OTPVerifiedEvent extends jspb.Message {
  getProfilename(): string;
  setProfilename(value: string): OTPVerifiedEvent;

  getIsnewuser(): boolean;
  setIsnewuser(value: boolean): OTPVerifiedEvent;

  getPhonenumber(): number;
  setPhonenumber(value: number): OTPVerifiedEvent;

  getLoginflow(): string;
  setLoginflow(value: string): OTPVerifiedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OTPVerifiedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: OTPVerifiedEvent): OTPVerifiedEvent.AsObject;
  static serializeBinaryToWriter(message: OTPVerifiedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OTPVerifiedEvent;
  static deserializeBinaryFromReader(message: OTPVerifiedEvent, reader: jspb.BinaryReader): OTPVerifiedEvent;
}

export namespace OTPVerifiedEvent {
  export type AsObject = {
    profilename: string,
    isnewuser: boolean,
    phonenumber: number,
    loginflow: string,
  }
}

export class AddToCartEvent extends jspb.Message {
  getProductid(): number;
  setProductid(value: number): AddToCartEvent;

  getProducttype(): string;
  setProducttype(value: string): AddToCartEvent;

  getVariantid(): number;
  setVariantid(value: number): AddToCartEvent;

  getQuantity(): number;
  setQuantity(value: number): AddToCartEvent;

  getChargeableamount(): number;
  setChargeableamount(value: number): AddToCartEvent;

  getCheckoutid(): string;
  setCheckoutid(value: string): AddToCartEvent;

  getStreamid(): number;
  setStreamid(value: number): AddToCartEvent;

  getIslive(): boolean;
  setIslive(value: boolean): AddToCartEvent;

  getHcpid(): number;
  setHcpid(value: number): AddToCartEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddToCartEvent.AsObject;
  static toObject(includeInstance: boolean, msg: AddToCartEvent): AddToCartEvent.AsObject;
  static serializeBinaryToWriter(message: AddToCartEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddToCartEvent;
  static deserializeBinaryFromReader(message: AddToCartEvent, reader: jspb.BinaryReader): AddToCartEvent;
}

export namespace AddToCartEvent {
  export type AsObject = {
    productid: number,
    producttype: string,
    variantid: number,
    quantity: number,
    chargeableamount: number,
    checkoutid: string,
    streamid: number,
    islive: boolean,
    hcpid: number,
  }
}

export class GoToPaymentsEvent extends jspb.Message {
  getProductid(): number;
  setProductid(value: number): GoToPaymentsEvent;

  getCheckoutid(): string;
  setCheckoutid(value: string): GoToPaymentsEvent;

  getProducttype(): string;
  setProducttype(value: string): GoToPaymentsEvent;

  getCoupon(): string;
  setCoupon(value: string): GoToPaymentsEvent;

  getChargeableamount(): number;
  setChargeableamount(value: number): GoToPaymentsEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoToPaymentsEvent.AsObject;
  static toObject(includeInstance: boolean, msg: GoToPaymentsEvent): GoToPaymentsEvent.AsObject;
  static serializeBinaryToWriter(message: GoToPaymentsEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoToPaymentsEvent;
  static deserializeBinaryFromReader(message: GoToPaymentsEvent, reader: jspb.BinaryReader): GoToPaymentsEvent;
}

export namespace GoToPaymentsEvent {
  export type AsObject = {
    productid: number,
    checkoutid: string,
    producttype: string,
    coupon: string,
    chargeableamount: number,
  }
}

export class InitiatePurchaseEvent extends jspb.Message {
  getIsbag(): boolean;
  setIsbag(value: boolean): InitiatePurchaseEvent;

  getProductid(): number;
  setProductid(value: number): InitiatePurchaseEvent;

  getProducttype(): string;
  setProducttype(value: string): InitiatePurchaseEvent;

  getCheckoutid(): string;
  setCheckoutid(value: string): InitiatePurchaseEvent;

  getChargeableamount(): number;
  setChargeableamount(value: number): InitiatePurchaseEvent;

  getCoupon(): string;
  setCoupon(value: string): InitiatePurchaseEvent;

  getPaymentmethod(): string;
  setPaymentmethod(value: string): InitiatePurchaseEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitiatePurchaseEvent.AsObject;
  static toObject(includeInstance: boolean, msg: InitiatePurchaseEvent): InitiatePurchaseEvent.AsObject;
  static serializeBinaryToWriter(message: InitiatePurchaseEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitiatePurchaseEvent;
  static deserializeBinaryFromReader(message: InitiatePurchaseEvent, reader: jspb.BinaryReader): InitiatePurchaseEvent;
}

export namespace InitiatePurchaseEvent {
  export type AsObject = {
    isbag: boolean,
    productid: number,
    producttype: string,
    checkoutid: string,
    chargeableamount: number,
    coupon: string,
    paymentmethod: string,
  }
}

export class CompletePurchaseEvent extends jspb.Message {
  getProductid(): number;
  setProductid(value: number): CompletePurchaseEvent;

  getProducttype(): string;
  setProducttype(value: string): CompletePurchaseEvent;

  getCheckoutid(): string;
  setCheckoutid(value: string): CompletePurchaseEvent;

  getChargeableamount(): number;
  setChargeableamount(value: number): CompletePurchaseEvent;

  getCoupon(): string;
  setCoupon(value: string): CompletePurchaseEvent;

  getPaymentmethod(): string;
  setPaymentmethod(value: string): CompletePurchaseEvent;

  getTransactionid(): string;
  setTransactionid(value: string): CompletePurchaseEvent;

  getProductidsList(): Array<number>;
  setProductidsList(value: Array<number>): CompletePurchaseEvent;
  clearProductidsList(): CompletePurchaseEvent;
  addProductids(value: number, index?: number): CompletePurchaseEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompletePurchaseEvent.AsObject;
  static toObject(includeInstance: boolean, msg: CompletePurchaseEvent): CompletePurchaseEvent.AsObject;
  static serializeBinaryToWriter(message: CompletePurchaseEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompletePurchaseEvent;
  static deserializeBinaryFromReader(message: CompletePurchaseEvent, reader: jspb.BinaryReader): CompletePurchaseEvent;
}

export namespace CompletePurchaseEvent {
  export type AsObject = {
    productid: number,
    producttype: string,
    checkoutid: string,
    chargeableamount: number,
    coupon: string,
    paymentmethod: string,
    transactionid: string,
    productidsList: Array<number>,
  }
}

export class PurchaseFailedEvent extends jspb.Message {
  getCheckoutid(): string;
  setCheckoutid(value: string): PurchaseFailedEvent;

  getChargeableamount(): number;
  setChargeableamount(value: number): PurchaseFailedEvent;

  getReason(): string;
  setReason(value: string): PurchaseFailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PurchaseFailedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: PurchaseFailedEvent): PurchaseFailedEvent.AsObject;
  static serializeBinaryToWriter(message: PurchaseFailedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PurchaseFailedEvent;
  static deserializeBinaryFromReader(message: PurchaseFailedEvent, reader: jspb.BinaryReader): PurchaseFailedEvent;
}

export namespace PurchaseFailedEvent {
  export type AsObject = {
    checkoutid: string,
    chargeableamount: number,
    reason: string,
  }
}

export class Response extends jspb.Message {
  getStatus(): string;
  setStatus(value: string): Response;

  getMessage(): string;
  setMessage(value: string): Response;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Response.AsObject;
  static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
  static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Response;
  static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
}

export namespace Response {
  export type AsObject = {
    status: string,
    message: string,
  }
}

