import React from 'react'
import Button from "react-bootstrap/Button";
import Jumbotron from "react-bootstrap/Jumbotron";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import {useHistory} from 'react-router-dom'
import "../App.css"


export const LandingPage = () => {
    const history = useHistory();
    const changeRoute = (path: string) => {
        history.replace({pathname: path})
    };

    return (
 
        <Jumbotron fluid bsPrefix="landing">
            <Container fluid>
                <Row>
                    <Col xs><Button variant="outline-light" className="float-right" style={{margin: "10px", textTransform: "uppercase"}} onClick={() => changeRoute('/login')}>Login</Button></Col>
                </Row> 
                <Row>
                    <Col xs><h1 className="landing-title">2ofClubs</h1></Col>
                </Row>
                <Row>
                    <Col xs> <h2 className="landing-subtitle">Find the community for you!</h2></Col>
                </Row>
                <Row >
                    <Col xs className="text-center">
                        <Button style={{fontSize:"1.25rem", backgroundColor:"white", color: "#696DE9", textTransform: "uppercase", padding: "1rem 2rem", margin: "10vh 0 0 0", borderRadius: "1rem"}} onClick={() => changeRoute('/Signup/signup')}>Sign up</Button>
                    </Col>
                </Row>
                <Row >
                    <Col xs className="text-center">
                        <h3 className="landing-text">CONTINUE WITHOUT ACCOUNT</h3>
                    </Col>
                </Row>
            </Container>
        </Jumbotron>

    ); 

};

