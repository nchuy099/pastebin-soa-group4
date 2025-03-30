package com.tony.pastecreate.service;

import com.tony.pastecreate.model.PasteEntity;
import com.tony.pastecreate.repository.PasteRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class PasteCreateService {
    private final PasteRepository pasteRepository;

    public PasteEntity createPaste(PasteEntity paste) {
        return pasteRepository.save(paste);
    }
}
