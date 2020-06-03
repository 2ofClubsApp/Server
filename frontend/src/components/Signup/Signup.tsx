import React from 'react'
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
const Signup = () => {
    return (
        <div>

            <Form>
                <Form.Group controlId={"Email"}>
                    <Form.Label>Email Address</Form.Label>
                    <Form.Control type={"email"} placeholder={"Email"}/>
                </Form.Group>
            </Form>

            <Form>
                <Form.Group controlId={"Password"}>
                    <Form.Label>Password</Form.Label>
                    <Form.Control type={"password"} placeholder={"Password"}/>
                </Form.Group>
            </Form>
            <Button onClick={signup}>Sign Up</Button>
        </div>
    )
};

const signup = () => {
    // axios.get("/students/Chris").then(response => {
    //     console.log(response)
    // }).catch(err => {
    //     console.log(err + "Unable to get student ;.;");
    // })
    console.log("Signing Up!")
};

export default Signup
