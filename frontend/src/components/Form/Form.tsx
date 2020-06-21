import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Button from "react-bootstrap/Button";
import React from "react";

type FormDefinition = {

    labels: any[]
    title: string
    buttonName: string
    formSubmit: any
}

const Form = (formInfo: FormDefinition) => {
    return (
        <Container className="form d-flex justify-content-center align-items-center" style={{width: "350px"}}>
            <Row>
                <Col xs lg="12">
                    <h2 className="text-center mb-4">{formInfo.title}</h2>
                    {formInfo.labels}
                    <div className="d-flex justify-content-center">
                        <Button className="align-middle mt-3 pl-4 pr-4 pt-2 pd-2 btn-grad"
                                onClick={formInfo.formSubmit}>
                            {formInfo.buttonName}
                        </Button>
                    </div>
                </Col>
            </Row>
        </Container>
    )
}

export default Form;
