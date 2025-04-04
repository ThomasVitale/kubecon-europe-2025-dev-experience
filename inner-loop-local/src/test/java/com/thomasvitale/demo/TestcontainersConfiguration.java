package com.thomasvitale.demo;

import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.testcontainers.containers.Network;
import org.testcontainers.utility.DockerImageName;

import io.github.microcks.testcontainers.MicrocksContainersEnsemble;

@TestConfiguration
public class TestcontainersConfiguration {

    private static Network network = Network.newNetwork();
    
    @Bean
    MicrocksContainersEnsemble microcksEnsemble(DynamicPropertyRegistry registry) {
        DockerImageName nativeImage = DockerImageName.parse("quay.io/microcks/microcks-uber:1.10.0-native")
            .asCompatibleSubstituteFor("quay.io/microcks/microcks-uber:1.9.0");
        
        MicrocksContainersEnsemble ensemble = new MicrocksContainersEnsemble(network, nativeImage)
            .withPostman()
            .withAsyncFeature()
            .withAccessToHost(true)
            .withMainArtifacts("friends-openapi.yml");
    
        registry.add("friends.base-url",
            () -> ensemble.getMicrocksContainer().getRestMockEndpoint("API Friends", "0.0.1"));

        return ensemble;
    }

}
