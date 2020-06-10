import React from 'react'
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Button from "react-bootstrap/Button";
import Formz from "./Form/Form";
import {useHistory} from 'react-router-dom';
import '../../app.css';

const Signup = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };
    const user = {
        controlId: "username",
        label: "Username",
        type: "text",
        placeholder: "Username"
    }
    const email = {
        controlId: "email",
        label: "Email Address",
        type: "email",
        placeholder: "Email"
    };
    const pass = {
        controlId: "password",
        label: "Password",
        type: "password",
        placeholder: "Password"
    };
    const passconfirm = {
        controlId: "confirmpassword",
        label: "Confirm Password",
        type: "password",
        placeholder: "Password"
    };
    return (
        <div> 
            <Button variant="outline-light" className="m-2 text-uppercase" onClick={() => changeRoute('/')}>Back to Home</Button>
            <Container className="form d-flex justify-content-center align-items-center w-25">
                <Row>
                    <Col>
                        <h2 className="text-center mb-4">Sign up</h2>
                        <Formz props={user}/>
                        <Formz props={email}/>
                        <Formz props={pass}/>
                        <Formz props={passconfirm}/>
                        <div className="d-flex justify-content-center">
                        <Button onClick={signup} className=" align-middle mt-3 pl-4 pr-4 pt-2 pd-2 btn-grad">Sign up</Button>
                        </div>
                    </Col>
                </Row>
            </Container>
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
