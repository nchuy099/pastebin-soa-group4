package com.tony.pastecreate.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.tony.pastecreate.model.PasteEntity;
import com.tony.pastecreate.repository.PasteRepository;

@Service
public class PasteCreateService {

    private final PasteRepository pasteRepository;

    // Constructor injection
    @Autowired
    public PasteCreateService(PasteRepository pasteRepository) {
        this.pasteRepository = pasteRepository;
    }

    // Your service methods here
    public PasteEntity createPaste(PasteEntity pasteEntity) {
        return pasteRepository.save(pasteEntity);
    }
}
