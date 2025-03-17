package com.nchuy099.pastebin_soa.service;

import com.nchuy099.pastebin_soa.dto.MonthlyStatsDTO;
import com.nchuy099.pastebin_soa.model.PasteEntity;

import java.util.Date;
import java.util.List;

public interface PasteService {
    PasteEntity createPaste(PasteEntity pasteData) throws Exception;

    PasteEntity getPasteById(String id) throws Exception;

    List<PasteEntity> getPublicPastes() throws Exception;

    MonthlyStatsDTO getMonthlyStats(Date month) throws Exception;
}
