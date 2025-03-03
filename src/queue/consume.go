package queue

func (pipe *Pipe) Consume(callback func([]byte, string)) error {

	msgs, err := pipe.Channel.Consume(pipe.Queue.Name, "", true, false, false, false, nil)
	if err != nil {
		pipe.Channel.Close()
		return err
	}

	go func() {
		for msg := range msgs {
			callback(msg.Body, msg.MessageId)
		}
	}()

	return nil

}
