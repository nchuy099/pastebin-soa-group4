package com.nchuy099.pastebin_soa.service.impl;

import com.nchuy099.pastebin_soa.dto.projection.MonthlyStatsProjection;

import com.nchuy099.pastebin_soa.dto.response.MonthlyStatsResponse;
import com.nchuy099.pastebin_soa.repository.PasteRepository;
import com.nchuy099.pastebin_soa.service.PasteService;

import lombok.RequiredArgsConstructor;

import org.springframework.stereotype.Service;

import java.time.YearMonth;
import java.util.NoSuchElementException;

@Service
@RequiredArgsConstructor
public class PasteServiceImpl implements PasteService {

    private final PasteRepository pasteRepository;

    @Override
    public MonthlyStatsResponse getMonthlyStats(YearMonth yearMonth) {
        int year = yearMonth.getYear();
        int month = yearMonth.getMonthValue();
        MonthlyStatsProjection projection = pasteRepository.getMonthlyStats(year, month)
                .orElseThrow(() -> new NoSuchElementException("No data for month: " + yearMonth));

        return new MonthlyStatsResponse(projection);
    }
}