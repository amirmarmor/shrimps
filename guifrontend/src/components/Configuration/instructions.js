import {Card} from "react-bootstrap"
import React from "react"

function Instructions() {
  return (
    <Card className="card-user">
      <div className="card-image">
        <img
          alt="..."
          src={
            require("assets/img/fish5.jpg")
              .default
          }
        />
      </div>
      <Card.Body>
        <div className="author">
          <img
            alt="..."
            className="avatar border-gray"
            src={require("assets/img/logo.png").default}
          />
          <h3 className="title">Instructions</h3>
        </div>
        <ul>
          <li><span style={{fontWeight: "bold"}}>Offset</span> - id of first camera, if there is a webcam on
            computer use 2 otherwise 0
          </li>
          <li><span style={{fontWeight: "bold"}}>Type</span> - Record Video or Images
            <ul>
              <li><span style={{fontWeight: "bold"}}>Video</span> - record video according to the rule schedule
              </li>
              <li><span style={{fontWeight: "bold"}}>Image</span> - record an image according to the rule schedule
              </li>
            </ul>
          </li>
          <li><span style={{fontWeight: "bold"}}>Recurring</span> - When to record video/image
            <ul>
              <li><span style={{fontWeight: "bold"}}>Daily</span> - record every day on the same hour (start) / for image see below
              </li>
              <li><span style={{fontWeight: "bold"}}>Hourly</span> - record every hour on the same minute (start) / for image see below
              </li>
              <li><span style={{fontWeight: "bold"}}>By Minute</span> - record every minute on same second (start) / for image see below
              </li>
            </ul>
          </li>
          <li><span style={{fontWeight: "bold"}}>Start</span> - specify the Hour / Minute / Second on which to start
            recording - disabled on image
          </li>
          <li><span style={{fontWeight: "bold"}}>Duration</span> - specify the length of recording in seconds
          </li>
          <li><span style={{fontWeight: "bold"}}>Duration</span> - in image node represents number of pictures to capture in time group - if daily is selected the interval between images will a day in milliseconds divided by this value, in if hour an hour in milliseconds divided by this vale etc.
          </li>
        </ul>
      </Card.Body>
    </Card>
  )
}

export default Instructions