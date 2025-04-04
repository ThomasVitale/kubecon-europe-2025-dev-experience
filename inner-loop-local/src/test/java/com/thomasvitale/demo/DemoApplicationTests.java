package com.thomasvitale.demo;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.context.TestPropertySource;
import org.springframework.test.web.reactive.server.WebTestClient;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@Import(TestcontainersConfiguration.class)
@TestPropertySource(properties = "spring.testcontainers.dynamic-property-registry-injection=allow")
class DemoApplicationTests {

	@Autowired
	WebTestClient webTestClient;

    @Test
	void friends() {
		webTestClient.get()
			.uri("/friends")
			.accept(MediaType.TEXT_PLAIN)
			.exchange()
			.expectStatus().isOk()
			.expectBodyList(String.class)
			.value(response -> assertThat(response).contains("Joey"));
	}

}
