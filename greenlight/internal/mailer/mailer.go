package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// 아래에서 이메일 템플릿을 저장하기 위해 embed.FS(임베디드 파일 시스템)
// 타입의 새 변수를 선언합니다. 그 위에 `//go:embed <path>` 형식의 주석 지시어가 있는데,
// 이는 ./templates 디렉터리의 내용을 templateFS 임베디드
// 파일 시스템 변수에 저장하고 싶다는 것을 Go에 알려줍니다.  ↓↓↓

//go:embed "templates"
var templateFS embed.FS

// 메일러 구조체를 정의하여 메일러 인스턴스(SMTP 서버에 연결하는 데 사용)와
// 이메일 발신자 정보(예: "Alice Smith <alice@example.com>")를 포함하도록 합니다.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	// 지정된 SMTP 서버 설정으로 새 mail.Dialer 인스턴스를 초기화합니다.
	// 또한 이메일을 보낼 때마다 5초의 시간 제한을 사용하도록 구성합니다.
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// 메일러 유형에 Send() 메서드를 정의합니다. 이 메서드는 수신자 이메일 주소를 첫 번째 매개변수로,
// 템플릿이 포함된 파일 이름을 두 번째 매개변수로, 템플릿에 대한 동적 데이터를 임의의 매개변수로 받습니다.
func (m Mailer) Send(recipient, templateFile string, data any) error {

	// ParseFS() 메서드를 사용하여 임베디드 파일 시스템에서 필요한 템플릿 파일을 구문 분석합니다.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// 명명된 템플릿 "subject"를 실행하여 동적 데이터를 전달하고 결과를 bytes.Buffer 변수에 저장합니다.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// 동일한 패턴에 따라 "plainBody" 템플릿을 실행하고 결과를 plainBody 변수에 저장합니다.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// 동일한 패턴에 따라 "plainBody" 템플릿을 실행하고 결과를 plainBody 변수에 저장합니다.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	/// mail.NewMessage() 함수를 사용하여 새 mail.Message 인스턴스를 초기화합니다.
	// 그런 다음 SetHeader() 메서드를 사용하여 이메일 수신자, 발신자 및 제목 헤더를 설정하고,
	// SetBody() 메서드를 사용하여 일반 텍스트 본문을 설정하고, AddAlternative() 메서드를 사용하여
	// HTML 본문을 설정합니다. 한 가지 주의할 점은 AddAlternative()는 항상 SetBody() 이후에 호출해야 한다는 것입니다.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// dialer에서 DialAndSend() 메서드를 호출하여 보낼 메시지를 전달합니다.
	// 그러면 SMTP 서버에 대한 연결이 열리고 메시지가 전송된 후 연결이 닫힙니다.
	// 시간 초과가 있는 경우 "다이얼 tcp: I/O 시간 초과" 오류를 반환합니다.
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
