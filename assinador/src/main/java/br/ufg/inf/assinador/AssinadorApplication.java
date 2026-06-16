package br.ufg.inf.assinador;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import br.ufg.inf.assinador.service.FakeSignatureService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.WebApplicationType;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.Arrays;

@SpringBootApplication
public class AssinadorApplication {

	public static void main(String[] args) {
		if (args.length > 0 && (args[0].equalsIgnoreCase("sign") || args[0].equalsIgnoreCase("validate"))) {
			SpringApplication app = new SpringApplication(AssinadorApplication.class);
			app.setWebApplicationType(WebApplicationType.NONE);
			executeLocalCommand(args);
		} else {
			SpringApplication.run(AssinadorApplication.class, args);
		}
	}

	private static void executeLocalCommand(String[] args) {
		String command = args[0];
		FakeSignatureService service = new FakeSignatureService();
		ObjectMapper mapper = new ObjectMapper();
		SignatureResponse response;

		try {
			if (command.equalsIgnoreCase("sign")) {
				SignRequest request = parseSignArgs(args);
				response = service.sign(request);
			} else {
				ValidateRequest request = parseValidateArgs(args);
				response = service.validate(request);
			}

			System.out.println(mapper.writeValueAsString(response));
			System.exit(response.isError() ? 1 : 0);

		} catch (Exception e) {
			SignatureResponse errorResponse = new SignatureResponse(null, false, e.getMessage(), true);
			try {
				System.out.println(mapper.writeValueAsString(errorResponse));
			} catch (Exception ignored) {}
			System.exit(1);
		}
	}

	private static SignRequest parseSignArgs(String[] args) {
		SignRequest request = new SignRequest();
		for (int i = 1; i < args.length; i++) {
			if (args[i].equals("--content") && i + 1 < args.length) request.setContent(args[++i]);
			if (args[i].equals("--token") && i + 1 < args.length) request.setToken(args[++i]);
		}
		return request;
	}

	private static ValidateRequest parseValidateArgs(String[] args) {
		ValidateRequest request = new ValidateRequest();
		for (int i = 1; i < args.length; i++) {
			if (args[i].equals("--content") && i + 1 < args.length) request.setContent(args[++i]);
			if (args[i].equals("--signature") && i + 1 < args.length) request.setSignature(args[++i]);
		}
		return request;
	}
}