package com.test.controller;

import com.oracle.stdio;
import com.test.controller.dto.TestObject;
import com.test.controller.dto.OtherObject;

@RestController
public class Controller {

private String someVar;

private int someWhat = 1;


    /*
        @DD:ENDPOINT 'api/testpoint'
        @DD:DESCRIPTION 'important testendpoint'
        @DD:PARAM int id 'id of user'
        @DD:PARAM string token 'security-token' @DD:NOTNULL
        @DD:PAYLOAD TestObject 'json object'
        @DD:TYPE post
        @DD:RESPONSE 200 json OtherObject
        @DD:RESPONSE 500 text
    */
    function testController(some input stuff){
        String test = "12334"
        return new OtherObject();
    }

    /*
        @DD:ENDPOINT 'api/testpoint/two'
        @DD:DESCRIPTION 'important testendpoint number two'
        @DD:PARAM int id 'id of user'
        @DD:PARAM string token 'security-token' @DD:NOTNULL
        @DD:PAYLOAD TestObject 'json object'
        @DD:TYPE post
        @DD:RESPONSE 200 json TestObject
        @DD:RESPONSE 500 text
    */
    function anotherTestController(some input stuff){
        String test = "12334"
        return new TestObject();
    }
}