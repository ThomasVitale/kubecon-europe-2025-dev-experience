package com.thomasvitale.demo;

import java.util.UUID;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;
import org.springframework.context.annotation.Bean;
import org.springframework.data.annotation.Id;
import org.springframework.data.repository.ListCrudRepository;
import org.springframework.web.client.RestClient;
import org.springframework.web.servlet.function.RouterFunction;
import org.springframework.web.servlet.function.RouterFunctions;
import org.springframework.web.servlet.function.ServerResponse;

@SpringBootApplication
@ConfigurationPropertiesScan
public class DemoApplication {

	public static void main(String[] args) {
		SpringApplication.run(DemoApplication.class, args);
	}

	@Bean
	RestClient friendRestClient(RestClient.Builder restClientBuilder, FriendsServiceProperties properties) {
		return restClientBuilder
			.baseUrl(properties.baseUrl())
			.build();
	}

	@Bean
	@SuppressWarnings("null")
	RouterFunction<ServerResponse> router(MovieRepository movieRepository, RestClient restClient) {
		return RouterFunctions.route()
			.GET("/", _ -> ServerResponse.ok().body("London, baby!"))
			.GET("/friends", _ -> ServerResponse.ok()
				.body(restClient.get().uri("/list")))
			.GET("/movies", _ -> ServerResponse.ok().body(movieRepository.findAll()))
			.POST("/movies", request -> ServerResponse.ok().body(
				movieRepository.save(request.body(Movie.class))
			))
			.build();
	}

}

@ConfigurationProperties(prefix = "friends")
record FriendsServiceProperties(String baseUrl) {}

interface MovieRepository extends ListCrudRepository<Movie,UUID> {}

record Movie(@Id UUID id, String name) {}
